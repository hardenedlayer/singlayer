package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type DocsResource struct {
	buffalo.Resource
}

func (v DocsResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	docs := &models.Docs{}
	var err error
	if c.Session().Get("is_admin").(bool) {
		err = tx.Order("category, subject, title").All(docs)
	} else {
		err = tx.Where("published=?", true).
			Order("category, subject, title").
			All(docs)
	}
	if err != nil {
		return err
	}
	c.Set("docs", docs)
	c.Set("categories", models.DocCategories())
	c.Set("subjects", models.DocSubjects())
	return c.Render(200, r.HTML("docs/index.html"))
}

func (v DocsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	doc := &models.Doc{}
	err := tx.Find(doc, c.Param("doc_id"))
	if err != nil {
		return err
	}
	c.Set("doc", doc)
	return c.Render(200, r.HTML("docs/show.html"))
}

func (v DocsResource) New(c buffalo.Context) error {
	c.Set("doc", &models.Doc{})
	return c.Render(200, r.HTML("docs/new.html"))
}

func (v DocsResource) Create(c buffalo.Context) error {
	doc := &models.Doc{}
	err := c.Bind(doc)
	if err != nil {
		return err
	}
	doc.SingleID = getCurrentSingle(c).ID
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(doc)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("doc", doc)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("docs/new.html"))
	}
	c.Flash().Add("success", "Doc was created successfully")
	return c.Redirect(302, "/docs/%s", doc.ID)
}

func (v DocsResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	doc := &models.Doc{}
	err := tx.Find(doc, c.Param("doc_id"))
	if err != nil {
		return err
	}
	c.Set("doc", doc)
	return c.Render(200, r.HTML("docs/edit.html"))
}

func (v DocsResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	doc := &models.Doc{}
	err := tx.Find(doc, c.Param("doc_id"))
	if err != nil {
		return err
	}
	err = c.Bind(doc)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(doc)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("doc", doc)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("docs/edit.html"))
	}
	c.Flash().Add("success", "Doc was updated successfully")
	return c.Redirect(302, "/docs/%s", doc.ID)
}

func (v DocsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	doc := &models.Doc{}
	err := tx.Find(doc, c.Param("doc_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(doc)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "Doc was destroyed successfully")
	return c.Redirect(302, "/docs")
}
