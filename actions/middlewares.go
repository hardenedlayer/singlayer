package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

func AuthenticateHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user_id := c.Session().Get("user_id")
		if user_id == nil {
			err := c.Redirect(http.StatusTemporaryRedirect, "/login")
			return err
		}
		err := next(c)
		return err
	}
}
