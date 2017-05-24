package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

// HomeHandler is a default handler to serve up
// a home page.
func MeHandler(c buffalo.Context) (err error) {
	tx := c.Value("tx").(*pop.Connection)
	single := &models.Single{}
	err = tx.Find(single, c.Session().Get("user_id"))
	if err != nil {
		return err
	}
	c.Logger().Debugf("------ %v ------", single)

	users := &models.Users{}
	err = tx.BelongsTo(single).All(users)
	if err != nil {
		return err
	}
	c.Logger().Debugf("------ %v ------", users)

	c.Set("users", users)

	c.Set("theme", "default")
	return c.Render(200, r.HTML("me.html"))
}
