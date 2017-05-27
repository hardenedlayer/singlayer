package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type SinglesResource struct {
	buffalo.Resource
}

func (v SinglesResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	singles := &models.Singles{}
	err := tx.All(singles)
	if err != nil {
		return err
	}
	c.Set("singles", singles)
	return c.Render(200, r.HTML("singles/index.html"))
}

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
		c.Set("theme", "admin")
		return c.Render(422, r.HTML("singles/edit.html"))
	}
	c.Flash().Add("success", "Single was updated successfully")
	return c.Redirect(302, "/singles/%s", single.ID)
}

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
