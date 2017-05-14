package actions

import (
	"net/http"
	"strings"

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

// call by every page requests
func SessionInfoHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user_id := c.Session().Get("user_id")
		if user_id != nil {
			c.Set("user_id", user_id)
			c.Set("user_name", c.Session().Get("user_name"))
			c.Set("user_mail", c.Session().Get("user_mail"))
			c.Set("user_icon", c.Session().Get("user_icon"))
			if str, ok := c.Session().Get("permissions").(string); ok {
				c.Set("user_is_admin", strings.Contains(str, ":admin:"))
			}
		}
		err := next(c)
		return err
	}
}
