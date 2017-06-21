package actions

import (
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type MailsResource struct {
	buffalo.Resource
}

// ADMIN PROTECTED
func (v MailsResource) List(c buffalo.Context) error {
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
	mails := &models.Mails{}
	q := tx.Paginate(page, pp)
	err = q.Order("created_at desc").All(mails)
	pager = q.Paginator
	if err != nil {
		return err
	}
	if len(*mails) == 0 && page > 1{
		return c.Redirect(302, "/mails")
	}

	c.Set("pager", pager)
	c.Set("mails", mails)
	return c.Render(200, r.HTML("mails/index.html"))
}

// ADMIN PROTECTED
func (v MailsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	mail := &models.Mail{}
	err := tx.Find(mail, c.Param("mail_id"))
	if err != nil {
		return err
	}
	c.Set("mail", mail)
	return c.Render(200, r.HTML("mails/show.html"))
}

// ADMIN PROTECTED
func (v MailsResource) New(c buffalo.Context) error {
	c.Set("mail", &models.Mail{})
	return c.Render(200, r.HTML("mails/new.html"))
}

// ADMIN PROTECTED
func (v MailsResource) Create(c buffalo.Context) error {
	mail := models.PrepareMail("", "", "")
	err := c.Bind(mail)
	if err != nil {
		return err
	}

	single := getCurrentSingle(c)
	err = mail.Send(single.ID)
	if err != nil {
		c.Logger().Errorf("CANNOT SEND A MAIL: %v", err)
	}
	c.Flash().Add("success", "Mail was sent successfully")
	return c.Redirect(302, "/mails/%s", mail.ID)
}

// ADMIN PROTECTED
func (v MailsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	mail := &models.Mail{}
	err := tx.Find(mail, c.Param("mail_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(mail)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "Mail was destroyed successfully")
	return c.Redirect(302, "/mails")
}
