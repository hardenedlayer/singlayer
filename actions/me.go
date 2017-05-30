package actions

import (
	"github.com/gobuffalo/buffalo"
)

func MeHandler(c buffalo.Context) (err error) {
	single := getCurrentSingle(c)
	users := single.Users()
	c.Logger().Debugf("MeHandler() got users: %v --", users)
	c.Set("users", users)
	return c.Render(200, r.HTML("me.html"))
}
