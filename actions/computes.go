package actions

import (
	"errors"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type ComputesResource struct {
	buffalo.Resource
}

func (v ComputesResource) List(c buffalo.Context) error {
	computes := &models.Computes{}

	pager := &pop.Paginator{}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 1
	}
	pp, err := strconv.Atoi(c.Param("pp"))
	if err != nil {
		pp = 20
	}
	if pp > 100 {
		pp = 100
	}

	if c.Session().Get("is_admin").(bool) {
		tx := c.Value("tx").(*pop.Connection)
		q := tx.Paginate(page, pp)
		err := q.Order("account_id, hostname").All(computes)
		pager = q.Paginator
		if err != nil {
			return err
		}
	} else {
		actor := c.Value("actor").(string)
		single := getCurrentSingle(c)
		if actor == "All" {
			c.Logger().Debugf("multi mode for single!")
			computes, pager = single.Computes(page, pp)
		} else {
			user := single.UserByUsername(c.Value("actor"))
			if user == nil {
				c.Logger().Errorf("SECURITY: cannot found user for %v", actor)
				c.Flash().Add("warning", "Oops! Who are you?")
			} else {
				c.Logger().Debugf("single mode for %v. full-sync.", actor)
				count, err := models.SyncComputes(user)
				if err != nil {
					return err
				}
				c.Logger().Debugf("%v new computes are synced.", count)
				computes, pager = user.Computes(page, pp)
			}
		}
	}
	if len(*computes) == 0 {
		if page > 1 {
			return c.Redirect(302, "/computes")
		}
		c.Flash().Add("warning", "No Compute Instances are found")
	}

	c.Set("pager", pager)
	c.Set("computes", computes)
	return c.Render(200, r.HTML("computes/index.html"))
}

func (v ComputesResource) Show(c buffalo.Context) error {
	compute, err := setCompute(c)
	if err != nil {
		return err
	}

	c.Set("compute", compute)
	return c.Render(200, r.HTML("computes/show.html"))
}

//// Custom functions:

// Find target user based on the context and permission. (SECURITY)
func setCompute(c buffalo.Context) (compute *models.Compute, err error) {
	tx := c.Value("tx").(*pop.Connection)
	compute = &models.Compute{}
	if c.Session().Get("is_admin").(bool) {
		err = tx.Find(compute, c.Param("compute_id"))
	} else {
		single := getCurrentSingle(c)
		compute = single.Compute(c.Param("compute_id"))
		if compute == nil {
			err = c.Error(404, errors.New("Compute Not Found"))
		}
	}
	return
}
