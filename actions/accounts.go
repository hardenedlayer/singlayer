package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

// AccountsResource is the resource for the account model
type AccountsResource struct {
	buffalo.Resource
}

// List gets all Accounts. This function is mapped to the the path
// GET /accounts
func (v AccountsResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	accounts := &models.Accounts{}
	err := tx.All(accounts)
	if err != nil {
		return err
	}
	c.Set("accounts", accounts)

	c.Set("theme", "admin")
	return c.Render(200, r.HTML("accounts/index.html"))
}

// Show gets the data for one Account. This function is mapped to
// the path GET /accounts/{account_id}
func (v AccountsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	account := &models.Account{}
	err := tx.Find(account, c.Param("account_id"))
	if err != nil {
		return err
	}
	c.Set("account", account)
	return c.Render(200, r.HTML("accounts/show.html"))
}

// Destroy deletes a account from the DB. This function is mapped
// to the path DELETE /accounts/{account_id}
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
