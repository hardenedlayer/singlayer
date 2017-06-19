package actions

import (
	"errors"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
	"github.com/satori/go.uuid"
)

type UsersResource struct {
	buffalo.Resource
}

// ADMIN PROTECTED
func (v UsersResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	users := &models.Users{}
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
		c.Set("user", user)
		c.Flash().Add("danger", "The user you submit is already exists.")
		return c.Render(422, r.HTML("users/new.html"))
	}

	err = user.Update()
	if err != nil {
		c.Set("user", user)
		c.Flash().Add("danger", "Cannot verify credential!")
		return c.Render(422, r.HTML("users/new.html"))
	}
	user.LastBatch = time.Now()

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
		return c.Render(422, r.HTML("users/new.html"))
	}
	updateActors(c, tx)
	c.Flash().Add("success", "User was created successfully")
	if b, o := c.Session().Get("is_admin").(bool); b && o {
		return c.Redirect(302, "/users")
	} else {
		return c.Redirect(302, "/me")
	}
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
	err = user.Update()
	if err != nil {
		c.Set("user", user)
		c.Flash().Add("danger", "Cannot update! Please check your credential.")
		return c.Render(422, r.HTML("users/edit.html"))
	}

	acc := &models.Account{ID: user.AccountId}
	err = acc.UpdateAndSave(user)
	if err != nil {
		c.Logger().Warnf("cannot save account: %v, %v", err, acc)
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
		return c.Render(422, r.HTML("users/edit.html"))
	}
	c.Flash().Add("success", "User was updated successfully")
	if b, o := c.Session().Get("is_admin").(bool); b && o {
		return c.Redirect(302, "/users")
	} else {
		return c.Redirect(302, "/me")
	}
}

func (v UsersResource) Destroy(c buffalo.Context) error {
	user, err := setUser(c)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	err = tx.Destroy(user)
	if err != nil {
		c.Logger().Errorf("cannot delete user: %v --", err)
		c.Flash().Add("danger", "Sorry, cannot delete user. try again later.")
		return c.Redirect(302, "/me")
	}
	updateActors(c, tx)
	c.Flash().Add("success", "User was destroyed successfully")
	if b, o := c.Session().Get("is_admin").(bool); b && o {
		return c.Redirect(302, "/users")
	} else {
		return c.Redirect(302, "/me")
	}
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

// update session variable actors after add/delete users.
func updateActors(c buffalo.Context, tx *pop.Connection) {
	var actors []string
	users := &models.Users{}
	_ = tx.Where("single_id=?", getCurrentSingle(c).ID).All(users)
	for _, u := range *users {
		actors = append(actors, u.Username)
	}
	actors = append(actors, "All")
	c.Session().Set("actors", actors)
	c.Session().Save()
	c.Logger().Infof("store actors %v into session.", actors)
}
