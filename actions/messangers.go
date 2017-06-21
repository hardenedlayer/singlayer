package actions

import (
	"errors"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

type MessangersResource struct {
	buffalo.Resource
}

// ADMIN PROTECTED
func (v MessangersResource) List(c buffalo.Context) error {
	pager := &pop.Paginator{}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pp, err := strconv.Atoi(c.Param("pp"))
	if err != nil || pp < 5 {
		pp = 20
	}
	if pp > 100 {
		pp = 100
	}

	messangers := &models.Messangers{}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Paginate(page, pp)
	err = q.Order("single_id, method").All(messangers)
	pager = q.Paginator
	if err != nil {
		return err
	}
	if len(*messangers) == 0 && page > 1 {
		return c.Redirect(302, "/messangers")
	}

	c.Set("pager", pager)
	c.Set("messangers", messangers)
	return c.Render(200, r.HTML("messangers/index.html"))
}

// ADMIN PROTECTED
func (v MessangersResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	messanger := &models.Messanger{}
	err := tx.Find(messanger, c.Param("messanger_id"))
	if err != nil {
		return err
	}
	c.Set("messanger", messanger)
	return c.Render(200, r.HTML("messangers/show.html"))
}

func (v MessangersResource) New(c buffalo.Context) error {
	c.Set("messanger", &models.Messanger{})
	return c.Render(200, r.HTML("messangers/new.html"))
}

func (v MessangersResource) Create(c buffalo.Context) error {
	messanger := &models.Messanger{}
	err := c.Bind(messanger)
	if err != nil {
		return err
	}
	messanger.SingleID = getCurrentSingle(c).ID

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(messanger)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("messanger", messanger)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("messangers/new.html"))
	}

	c.Flash().Add("success", "Messanger was created successfully")
	if b, o := c.Session().Get("is_admin").(bool); b && o {
		return c.Redirect(302, "/messangers")
	} else {
		return c.Redirect(302, "/me")
	}
}

func (v MessangersResource) Edit(c buffalo.Context) error {
	messanger, err := setMessanger(c)
	if err != nil {
		return err
	}
	c.Set("messanger", messanger)
	return c.Render(200, r.HTML("messangers/edit.html"))
}

func (v MessangersResource) Update(c buffalo.Context) error {
	messanger, err := setMessanger(c)
	if err != nil {
		return err
	}
	err = c.Bind(messanger)
	if err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndUpdate(messanger)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("messanger", messanger)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("messangers/edit.html"))
	}

	c.Flash().Add("success", "Messanger was updated successfully")
	if b, o := c.Session().Get("is_admin").(bool); b && o {
		return c.Redirect(302, "/messangers")
	} else {
		return c.Redirect(302, "/me")
	}
}

func (v MessangersResource) Destroy(c buffalo.Context) error {
	messanger, err := setMessanger(c)
	if err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	err = tx.Destroy(messanger)
	if err != nil {
		return err
	}

	c.Flash().Add("success", "Messanger was destroyed successfully")
	if b, o := c.Session().Get("is_admin").(bool); b && o {
		return c.Redirect(302, "/messangers")
	} else {
		return c.Redirect(302, "/me")
	}
}

//// action helper

// protected set messanger
func setMessanger(c buffalo.Context) (m *models.Messanger, err error) {
	tx := c.Value("tx").(*pop.Connection)
	m = &models.Messanger{}
	if c.Session().Get("is_admin").(bool) {
		err = tx.Find(m, c.Param("messanger_id"))
	} else {
		single := getCurrentSingle(c)
		m = single.Messanger(c.Param("messanger_id"))
		if m == nil {
			err = c.Error(404, errors.New("Messanger Not Found"))
		}
	}
	c.Logger().Debugf("setMessanger() returns messanger: %v, err: %v", m, err)
	return
}
