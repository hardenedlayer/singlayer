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

func (v AccountsResource) Edit(c buffalo.Context) error {
	account, err := setAccount(c)
	if err != nil {
		return err
	}
	c.Set("account", account)
	return c.Render(200, r.HTML("accounts/edit.html"))
}

func (v AccountsResource) Update(c buffalo.Context) error {
	account, err := setAccount(c)
	if err != nil {
		return err
	}
	err = c.Bind(account)
	if err != nil {
		c.Logger().Errorf("cannot bind with new data: %v --", err)
		return err
	}

	user := getCurrentSingle(c).UserByAccount(c.Param("account_id"))
	err = account.UpdateAndSave(user)
	if err != nil {
		c.Logger().Warnf("cannot save account: %v, %v", err, account)
	}
	c.Flash().Add("success", "Account was updated successfully")
	return c.Redirect(302, "/me")
}

// Find target account based on the context and permission.
func setAccount(c buffalo.Context) (account *models.Account, err error) {
	tx := c.Value("tx").(*pop.Connection)
	account = &models.Account{}
	if c.Session().Get("is_admin").(bool) {
		err = tx.Find(account, c.Param("account_id"))
	} else {
		single := getCurrentSingle(c)
		account = single.Account(c.Param("account_id"))
		if account == nil {
			err = c.Error(404, errors.New("Account Not Found"))
		}
	}
	c.Logger().Debugf("setAccount() returns account: %v, err: %v", account, err)
	return
}
