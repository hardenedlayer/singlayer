package actions

import (
	"github.com/gobuffalo/buffalo"

	"github.com/hardenedlayer/singlayer/models"
)

func MeHandler(c buffalo.Context) (err error) {
	single := getCurrentSingle(c)

	c.Set("single", single)
	c.Set("users", single.Users())
	c.Set("user", &models.User{})
	c.Set("messangers", single.Messangers())
	c.Set("messanger", &models.Messanger{})
	return c.Render(200, r.HTML("me.html"))
}
