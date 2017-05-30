package actions

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gorilla/sessions"

	"github.com/hardenedlayer/singlayer/models"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"

	"github.com/markbates/goth/gothic"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:          ENV,
			SessionStore: newSessionStore(ENV),
			SessionName:  "_singlayer_session",
		})
		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(middleware.CSRF)

		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		T, err = i18n.New(packr.NewBox("../locales"), "en-US")
		if err != nil {
			log.Fatal(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)
		app.Middleware.Skip(AuthenticateHandler, HomeHandler)

		app.ServeFiles("/assets", packr.NewBox("../public/assets"))
		// route for authentication
		auth := app.Group("/auth")
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)

		app.GET("/login", LoginHandler)
		app.GET("/logout", LogoutHandler)
		app.Middleware.Skip(AuthenticateHandler, LoginHandler, LogoutHandler)

		app.Use(AuthenticateHandler)
		app.Use(SessionInfoHandler)

		// special routes without resource
		app.GET("/me", MeHandler)
		n := app.Group("/n")
		n.Use(AdminPageKeeper)
		n.GET("/meta/tickets", TicketsMetaHandler)
		n.GET("/ticket_statuses/sync", SyncTicketStatuses)
		n.GET("/ticket_subjects/sync", SyncTicketSubjects)
		n.GET("/ticket_groups/sync", SyncTicketGroups)

		// resource based routes
		var r buffalo.Resource

		g := app.Resource("/singles", SinglesResource{&buffalo.BaseResource{}})
		g.Use(AdminPageKeeper)

		r = &UsersResource{&buffalo.BaseResource{}}
		g = app.Resource("/users", r)
		g.Use(AdminPageKeeper)
		g.Middleware.Skip(AdminPageKeeper,
			r.Show, r.New, r.Edit, r.Create, r.Update, r.Destroy)

		r = &AccountsResource{&buffalo.BaseResource{}}
		g = app.Resource("/accounts", r)
		g.Use(AdminPageKeeper)
		g.Middleware.Skip(AdminPageKeeper, r.Show)

		r = &TicketsResource{&buffalo.BaseResource{}}
		g = app.Resource("/tickets", r)
		g.Use(AdminPageKeeper)
		g.Middleware.Skip(AdminPageKeeper, r.List, r.Show)

		r = &DirectLinksResource{&buffalo.BaseResource{}}
		g = app.Resource("/directlinks", r)
		g.Use(AdminPageKeeper)
		g.Middleware.Skip(AdminPageKeeper,
			r.List, r.Show, r.New, r.Create, r.Edit, r.Update)
	}

	return app
}

func newSessionStore(env string) sessions.Store {
	secret := envy.Get("SESSION_SECRET", "")
	if env == "production" && secret == "" {
		log.Fatal("FATAL! set SESSION_SECRET env variable for your security!")
	}
	cookieStore := sessions.NewCookieStore([]byte(secret))
	cookieStore.MaxAge(60 * 60 * 1)
	return cookieStore
}
