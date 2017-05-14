package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func AuthenticateHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user_id := c.Session().Get("user_id")
		if user_id == nil {
			c.Flash().Add("danger", "Sorry, Login required.")
			err := c.Redirect(http.StatusTemporaryRedirect, "/login")
			return err
		}
		err := next(c)
		return err
	}
}

func AdminPageKeeper(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		is_admin := c.Session().Get("is_admin")
		if is_admin == false {
			c.Flash().Add("danger", "STAFF ONLY")
			err := c.Redirect(http.StatusTemporaryRedirect, "/")
			return err
		}
		err := next(c)
		return err
	}
}

// call by every page requests
func SessionInfoHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user_id := c.Session().Get("user_id")
		if user_id != nil {
			c.Set("user_id", user_id)
			c.Set("user_name", c.Session().Get("user_name"))
			c.Set("user_mail", c.Session().Get("user_mail"))
			c.Set("user_icon", c.Session().Get("user_icon"))
			c.Set("user_is_admin", c.Session().Get("is_admin"))
		}
		err := next(c)
		return err
	}
}
