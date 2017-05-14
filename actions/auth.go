package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/gplus/callback")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/facebook/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	session := c.Session()
	session.Set("user_id", user.UserID)
	session.Set("user_name", user.Name)
	session.Set("user_mail", user.Email)
	session.Set("user_icon", user.AvatarURL)
	session.Set("permissions", ":guest:")
	err = session.Save()
	if err != nil {
		return c.Error(401, err)
	}
	c.Flash().Add("success", "You have been successfully logged in.")
	return c.Redirect(307, "/")
}

func LoginHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("login.html"))
}

func LogoutHandler(c buffalo.Context) error {
	session := c.Session()
	session.Clear()
	session.Save()
	c.Flash().Add("success", "You have been successfully logged out.")
	return c.Redirect(307, "/")
}
