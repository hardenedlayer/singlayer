package actions

import (
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

type SinglesResource struct {
	buffalo.Resource
}

// ADMIN PROTECTED
func (v SinglesResource) List(c buffalo.Context) error {
	singles := &models.Singles{}
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

	tx := c.Value("tx").(*pop.Connection)
	q := tx.Paginate(page, pp)
	err = q.Order("permissions").All(singles)
	pager = q.Paginator
	if err != nil {
		return err
	}
	if len(*singles) == 0 && page > 1{
		return c.Redirect(302, "/singles")
	}

	c.Set("pager", pager)
	c.Set("singles", singles)
	return c.Render(200, r.HTML("singles/index.html"))
}

// ADMIN PROTECTED
func (v SinglesResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	single := &models.Single{}
	err := tx.Find(single, c.Param("single_id"))
	if err != nil {
		return err
	}
	c.Set("single", single)
	return c.Render(200, r.HTML("singles/show.html"))
}

// ADMIN PROTECTED
func (v SinglesResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	single := &models.Single{}
	err := tx.Find(single, c.Param("single_id"))
	if err != nil {
		return err
	}
	c.Set("single", single)
	return c.Render(200, r.HTML("singles/edit.html"))
}

// ADMIN PROTECTED
func (v SinglesResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	single := &models.Single{}
	err := tx.Find(single, c.Param("single_id"))
	if err != nil {
		return err
	}
	err = c.Bind(single)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(single)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("single", single)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("singles/edit.html"))
	}
	c.Flash().Add("success", "Single was updated successfully")
	return c.Redirect(302, "/singles/%s", single.ID)
}

// ADMIN PROTECTED
func (v SinglesResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	single := &models.Single{}
	err := tx.Find(single, c.Param("single_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(single)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "Single was destroyed successfully")
	return c.Redirect(302, "/singles")
}
