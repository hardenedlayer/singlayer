package actions

import (
	"time"
	"errors"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/jinzhu/copier"
	"github.com/markbates/pop"
	"github.com/satori/go.uuid"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type UsersResource struct {
	buffalo.Resource
}

// ADMIN PROTECTED
func (v UsersResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	users := &models.Users{}
	// TODO: order, paging,...
	err := tx.Order("username asc").All(users)
	if err != nil {
		return err
	}
	c.Set("users", users)
	return c.Render(200, r.HTML("users/index.html"))
}

func (v UsersResource) Show(c buffalo.Context) error {
	user, err := setUser(c)
	if err != nil {
		return err
	}
	c.Set("user", user)
	return c.Render(200, r.HTML("users/show.html"))
}

func (v UsersResource) New(c buffalo.Context) error {
	c.Set("user", &models.User{})
	return c.Render(200, r.HTML("users/new.html"))
}

func (v UsersResource) Edit(c buffalo.Context) error {
	user, err := setUser(c)
	if err != nil {
		return err
	}
	c.Set("user", user)
	return c.Render(200, r.HTML("users/edit.html"))
}

func (v UsersResource) Create(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	err := c.Bind(user)
	if err != nil {
		return err
	}
	if id, ok := c.Session().Get("user_id").(uuid.UUID); ok {
		user.SingleID = id
	}

	users := &models.Users{}
	err = tx.Where("username=?", user.Username).All(users)
	if err != nil {
		return c.Error(501, err)
	}
	if len(*users) != 0 {
		c.Flash().Add("danger", "The user you submit is already exists.")
		return c.Redirect(302, "/users/new")
	}

	err = setupUser(c, user)
	if err != nil {
		c.Set("user", user)
		c.Set("error", err)
		c.Logger().Errorf("SETUP ERROR: %v --", err)
		c.Set("theme", "admin")
		return c.Render(422, r.HTML("users/edit.html"))
	}

	acc := &models.Account{ID: user.AccountId}
	err = acc.UpdateAndSave(user)
	if err != nil {
		c.Logger().Warnf("cannot save account: %v, %v", err, acc)
	}

	c.Logger().Debugf("about to create an user: %v ----", user)
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		c.Logger().Errorf("validation error: %v --", err)
		return err
	}
	if verrs.HasAny() {
		c.Set("user", user)
		c.Set("errors", verrs)
		c.Logger().Errorf("validation errors: %v --", verrs)
		c.Set("theme", "admin")
		return c.Render(422, r.HTML("users/new.html"))
	}
	c.Flash().Add("success", "User was created successfully")
	return c.Redirect(302, "/users/%d", user.ID)
}

func (v UsersResource) Update(c buffalo.Context) error {
	user, err := setUser(c)
	if err != nil {
		return err
	}
	err = c.Bind(user)
	if err != nil {
		c.Logger().Errorf("cannot bind with new data: %v --", err)
		return err
	}
	err = setupUser(c, user)
	if err != nil {
		c.Set("user", user)
		c.Set("error", err)
		c.Logger().Errorf("SETUP ERROR: %v --", err)
		c.Set("theme", "admin")
		return c.Render(422, r.HTML("users/edit.html"))
	}

	c.Logger().Debugf("about to update an user: %v ----", user)
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		c.Logger().Errorf("validation error: %v --", err)
		return err
	}
	if verrs.HasAny() {
		c.Set("user", user)
		c.Set("errors", verrs)
		c.Logger().Errorf("validation errors: %v --", verrs)
		c.Set("theme", "admin")
		return c.Render(422, r.HTML("users/edit.html"))
	}
	c.Flash().Add("success", "User was updated successfully")
	return c.Redirect(302, "/users/%d", user.ID)
}

func (v UsersResource) Destroy(c buffalo.Context) error {
	user, err := setUser(c)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	err = tx.Destroy(user)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "User was destroyed successfully")
	return c.Redirect(302, "/me")
}

// Find target user based on the context and permission.
func setUser(c buffalo.Context) (user *models.User, err error) {
	tx := c.Value("tx").(*pop.Connection)
	user = &models.User{}
	if c.Session().Get("is_admin").(bool) {
		err = tx.Find(user, c.Param("user_id"))
	} else {
		single := getCurrentSingle(c)
		user = single.User(c.Param("user_id"))
		if user == nil {
			err = c.Error(404, errors.New("User Not Found"))
		}
	}
	c.Logger().Debugf("setUser() returns user: %v, err: %v", user, err)
	return
}

// Fill up user struct with response of api call.
func setupUser(c buffalo.Context, user *models.User) (err error) {
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"
	service := services.GetAccountService(sess)
	sl_user, err := service.
		Mask("id;accountId;parentId;companyName;email;firstName;lastName;ticketCount;openTicketCount;hardwareCount;virtualGuestCount").
		GetCurrentUser()
	if err != nil {
		c.Logger().Errorf("softlayer api exception: %v --", err)
		return err
	}
	copier.Copy(user, sl_user)
	user.ID = *sl_user.Id
	user.LastBatch = time.Now()
	c.Logger().Debugf("check user: %v ----", user)
	return
}
