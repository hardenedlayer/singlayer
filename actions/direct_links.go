package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type DirectLinksResource struct {
	buffalo.Resource
}

func (v DirectLinksResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	dlinks := &models.DirectLinks{}
	err := tx.All(dlinks)
	if err != nil {
		return err
	}
	c.Set("dlinks", dlinks)
	return c.Render(200, r.HTML("direct_links/index.html"))
}

func (v DirectLinksResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	dlink := &models.DirectLink{}
	err := tx.Find(dlink, c.Param("directlink_id"))
	if err != nil {
		return err
	}
	c.Set("dlink", dlink)
	return c.Render(200, r.HTML("direct_links/show.html"))
}

func (v DirectLinksResource) New(c buffalo.Context) error {
	dlink := &models.DirectLink{}
	dlink.VlanId = 999
	dlink.Type = "CX"
	dlink.Location = "SEO01"
	dlink.LineNumber = 1
	dlink.Port = "N/A"
	dlink.Router = "#1"
	dlink.Speed = 1
	dlink.RoutingOption = "Local"
	dlink.Prefix = 31
	dlink.Migration = "ANY"
	dlink.Signature = "ANY"
	dlink.Status = "ordered"
	c.Set("dlink", dlink)
	return c.Render(200, r.HTML("direct_links/new.html"))
}

func (v DirectLinksResource) Create(c buffalo.Context) error {
	dlink := &models.DirectLink{}
	err := c.Bind(dlink)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(dlink)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("dlink", dlink)
		c.Set("errors", verrs)
		c.Logger().Printf("errors %v", verrs)
		return c.Render(422, r.HTML("direct_links/new.html"))
	}
	c.Flash().Add("success", "DirectLink was created successfully")
	return c.Redirect(302, "/directlinks/%s", dlink.ID)
}

func (v DirectLinksResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	dlink := &models.DirectLink{}
	err := tx.Find(dlink, c.Param("directlink_id"))
	if err != nil {
		return err
	}
	c.Set("dlink", dlink)
	return c.Render(200, r.HTML("direct_links/edit.html"))
}

func (v DirectLinksResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	dlink := &models.DirectLink{}
	err := tx.Find(dlink, c.Param("directlink_id"))
	if err != nil {
		return err
	}
	err = c.Bind(dlink)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(dlink)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("dlink", dlink)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("direct_links/edit.html"))
	}
	c.Flash().Add("success", "DirectLink was updated successfully")
	return c.Redirect(302, "/directlinks/%s", dlink.ID)
}

// ADMIN PROTECTED
func (v DirectLinksResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	dlink := &models.DirectLink{}
	err := tx.Find(dlink, c.Param("directlink_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(dlink)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "DirectLink was destroyed successfully")
	return c.Redirect(302, "/directlinks")
}
