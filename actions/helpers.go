package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

func getCurrentSingle(c buffalo.Context) (single *models.Single) {
	tx := c.Value("tx").(*pop.Connection)
	single = &models.Single{}
	err := tx.Find(single, c.Session().Get("user_id"))
	if err != nil {
		return nil
	}
	c.Logger().Debugf("getCurrentSingle() returns single: %v", single)
	return single
}
