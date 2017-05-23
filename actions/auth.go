package actions

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/gplus/callback")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/facebook/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	c.Logger().Debugf("user: %v ---", user)

	singles := &models.Singles{}
	single := &models.Single{}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("provider=?", user.Provider).Where("user_id=?", user.UserID)

	err = q.All(singles)
	if err != nil {
		// TODO add error recognition code and check the error page:
		return c.Error(501, err)
	}

	if len(*singles) == 1 {
		// TODO action logging
		single = &(*singles)[0]
		c.Flash().Add("success", "Welcome back! I missed you...")
	} else if len(*singles) == 0 {
		if user.Email == "" {
			c.Flash().Add("danger",
				"Sorry but unacceptable account. (no email provided)")
			return c.Redirect(307, "/login")
		}
		// TODO action logging
		single.Provider = user.Provider
		single.Email = user.Email
		single.Name = user.Name
		single.UserID = user.UserID
		single.AvatarUrl = user.AvatarURL
		single.Permissions = ":guest:"
		// TODO mark as admin for v ery first user.
		verrs, err := tx.ValidateAndCreate(single)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return c.Error(501, err)
		}
		if verrs.HasAny() {
			fmt.Printf("verrs: %v\n", verrs)
			return c.Error(501, verrs)
		}

		err = q.First(single)
		c.Flash().Add("info", "Nice to meet you! You just become a singler!")
	} else {
		// TODO action logging
		return c.Error(501, errors.New("Somthing went wrong!!!"))
	}
	fmt.Printf("final single: %v\n", single)

	session := c.Session()
	session.Set("user_id", single.ID)
	session.Set("user_name", single.Name)
	session.Set("user_mail", single.Email)
	session.Set("user_icon", single.AvatarUrl)
	session.Set("permissions", single.Permissions)
	session.Set("is_admin", strings.Contains(single.Permissions, ":admin:"))
	err = session.Save()
	if err != nil {
		return c.Error(401, err)
	}
	return c.Redirect(307, "/")
}

func LoginHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("login.html"))
}

func LogoutHandler(c buffalo.Context) error {
	session := c.Session()
	session.Clear()
	session.Save()
	c.Flash().Add("success", "You have been successfully logged out.")
	return c.Redirect(307, "/")
}
