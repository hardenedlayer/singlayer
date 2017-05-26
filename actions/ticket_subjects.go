package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/hardenedlayer/singlayer/models"
	"github.com/markbates/pop"
)

type TicketSubjectsResource struct {
	buffalo.Resource
}

// Sync all TicketSubjects from original site.
func SyncTicketSubjects(c buffalo.Context) error {
	user := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Where("single_id=?", c.Session().Get("user_id")).First(user)
	if err != nil {
		c.Logger().Errorf("TX error: %v", err)
	}
	c.Logger().Infof("using %v:%v", user.Username, user.APIKey)
	err = models.SyncTicketSubjects(user)
	if err == nil {
		c.Flash().Add("success", "TicketSubjects were synced successfully")
	} else {
		c.Flash().Add("danger", "Cannot sync TicketSubjects")
	}
	return c.Redirect(302, "/n/meta/tickets")
}

// List gets all TicketSubjects. This function is mapped to the the path
// GET /ticket_subjects
func (v TicketSubjectsResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketSubjects := &models.TicketSubjects{}
	err := tx.All(ticketSubjects)
	if err != nil {
		return err
	}
	c.Set("ticketSubjects", ticketSubjects)
	return c.Render(200, r.HTML("ticket_subjects/index.html"))
}

// Show gets the data for one TicketSubject. This function is mapped to
// the path GET /ticket_subjects/{ticket_subject_id}
func (v TicketSubjectsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketSubject := &models.TicketSubject{}
	err := tx.Find(ticketSubject, c.Param("ticket_subject_id"))
	if err != nil {
		return err
	}
	c.Set("ticketSubject", ticketSubject)
	return c.Render(200, r.HTML("ticket_subjects/show.html"))
}

// New renders the formular for creating a new ticket_subject.
// This function is mapped to the path GET /ticket_subjects/new
func (v TicketSubjectsResource) New(c buffalo.Context) error {
	c.Set("ticketSubject", &models.TicketSubject{})
	return c.Render(200, r.HTML("ticket_subjects/new.html"))
}

// Create adds a ticket_subject to the DB. This function is mapped to the
// path POST /ticket_subjects
func (v TicketSubjectsResource) Create(c buffalo.Context) error {
	ticketSubject := &models.TicketSubject{}
	err := c.Bind(ticketSubject)
	if err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(ticketSubject)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticketSubject", ticketSubject)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("ticket_subjects/new.html"))
	}
	c.Flash().Add("success", "TicketSubject was created successfully")
	return c.Redirect(302, "/ticket_subjects/%d", ticketSubject.ID)
}

// Edit renders a edit formular for a ticket_subject. This function is
// mapped to the path GET /ticket_subjects/{ticket_subject_id}/edit
func (v TicketSubjectsResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketSubject := &models.TicketSubject{}
	err := tx.Find(ticketSubject, c.Param("ticket_subject_id"))
	if err != nil {
		return err
	}
	c.Set("ticketSubject", ticketSubject)
	return c.Render(200, r.HTML("ticket_subjects/edit.html"))
}

// Update changes a ticket_subject in the DB. This function is mapped to
// the path PUT /ticket_subjects/{ticket_subject_id}
func (v TicketSubjectsResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketSubject := &models.TicketSubject{}
	err := tx.Find(ticketSubject, c.Param("ticket_subject_id"))
	if err != nil {
		return err
	}
	err = c.Bind(ticketSubject)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(ticketSubject)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("ticketSubject", ticketSubject)
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("ticket_subjects/edit.html"))
	}
	c.Flash().Add("success", "TicketSubject was updated successfully")
	return c.Redirect(302, "/ticket_subjects/%d", ticketSubject.ID)
}

// Destroy deletes a ticket_subject from the DB. This function is mapped
// to the path DELETE /ticket_subjects/{ticket_subject_id}
func (v TicketSubjectsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	ticketSubject := &models.TicketSubject{}
	err := tx.Find(ticketSubject, c.Param("ticket_subject_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(ticketSubject)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "TicketSubject was destroyed successfully")
	return c.Redirect(302, "/ticket_subjects")
}
