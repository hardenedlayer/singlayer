package actions

import (
	"errors"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type AccountsResource struct {
	buffalo.Resource
}

// ADMIN PROTECTED
func (v AccountsResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	accounts := &models.Accounts{}
	err := tx.All(accounts)
	if err != nil {
		return err
	}
	c.Set("accounts", accounts)
	return c.Render(200, r.HTML("accounts/index.html"))
}

func (v AccountsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	account := &models.Account{}
	if c.Session().Get("is_admin").(bool) {
		err := tx.Find(account, c.Param("account_id"))
		if err != nil {
			return err
		}
	} else {
		single := getCurrentSingle(c)
		account = single.Account(c.Param("account_id"))
		c.Logger().Debugf("single: %v", account)
		if account == nil {
			return c.Error(404, errors.New("Account Not Found"))
		}
	}
	c.Set("account", account)
	return c.Render(200, r.HTML("accounts/show.html"))
}

// ADMIN PROTECTED
func (v AccountsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	account := &models.Account{}
	err := tx.Find(account, c.Param("account_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(account)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "Account was destroyed successfully")
	return c.Redirect(302, "/accounts")
}
