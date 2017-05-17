package actions

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
	"github.com/satori/go.uuid"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/services"
)

// UsersResource is the resource for the user model
type UsersResource struct {
	buffalo.Resource
}


// Find a single user or list of users to show
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
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		return err
	}
	c.Set("user", user)
	return c.Render(200, r.HTML("users/show.html"))
}


// New and Edit: generated a form for create and update
func (v UsersResource) New(c buffalo.Context) error {
	c.Set("user", &models.User{})
	return c.Render(200, r.HTML("users/new.html"))
}

func (v UsersResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		c.Logger().Errorf("cannot found user to edit: %v --", err)
		return err
	}
	c.Set("user", user)
	return c.Render(200, r.HTML("users/edit.html"))
}


// Create and Update: serve user's request for insert and update
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
	setupUser(c, user)

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
		return c.Render(422, r.HTML("users/new.html"))
	}
	c.Flash().Add("success", "User was created successfully")
	return c.Redirect(302, "/users/%d", user.ID)
}

func (v UsersResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		c.Logger().Errorf("cannot found user to update: %v --", err)
		return err
	}
	err = c.Bind(user)
	if err != nil {
		c.Logger().Errorf("cannot bind with new data: %v --", err)
		return err
	}

	c.Logger().Debugf("about to update an user: %v ----", user)
	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		c.Logger().Errorf("validation error: %v --", err)
		return err
	}
	if verrs.HasAny() {
		c.Set("user", user)
		c.Set("errors", verrs)
		c.Logger().Errorf("validation errors: %v --", verrs)
		return c.Render(422, r.HTML("users/edit.html"))
	}
	c.Flash().Add("success", "User was updated successfully")
	return c.Redirect(302, "/users/%d", user.ID)
}

// Delete
func (v UsersResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(user)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "User was destroyed successfully")
	return c.Redirect(302, "/users")
}



func setupUser(c buffalo.Context, user *models.User) (err error) {
	sess := session.New(user.Username, user.APIKey)
	sess.Endpoint = "https://api.softlayer.com/rest/v3.1"
	service := services.GetAccountService(sess)
	sl_user, err := service.GetCurrentUser()
	if err != nil {
		c.Logger().Errorf("cannot get current user from session: %v --", err)
		return err
	}
	user.ID = *sl_user.Id
	user.AccountId = *sl_user.AccountId
	user.ParentId = *sl_user.ParentId
	user.CompanyName = *sl_user.CompanyName
	user.Email = *sl_user.Email
	user.FirstName = *sl_user.FirstName
	user.LastName = *sl_user.LastName
	user.LastBatch = time.Now()
	c.Logger().Debugf("check user: %v ----", user)
	return
}
