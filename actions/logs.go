package actions

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"

	"github.com/hardenedlayer/singlayer/models"
)

type LogsResource struct {
	buffalo.Resource
}

func (v LogsResource) List(c buffalo.Context) error {
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
	logs := &models.Logs{}
	q := tx.Paginate(page, pp).Where("is_read = ?", false)
	err = q.Order("created_at desc").All(logs)
	pager = q.Paginator
	if err != nil {
		return err
	}
	if len(*logs) == 0 && page > 1 {
		return c.Redirect(302, "/logss")
	}

	c.Set("pager", pager)
	c.Set("logs", logs)
	return c.Render(200, r.HTML("logs/index.html"))
}

func (v LogsResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	log := &models.Log{}
	err := tx.Find(log, c.Param("log_id"))
	if err != nil {
		return err
	}

	log.IsRead = true
	err = tx.Update(log)
	if err != nil {
		return err
	}
	c.Flash().Add("success", "Log was marked as read")
	return c.Redirect(302, "/logs")
}

////

func l(c buffalo.Context, cat, lev, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)

	switch lev {
	case INFO:
		c.Logger().Info("APPLOG: " + message)
	case WARN:
		c.Logger().Warn("APPLOG: " + message)
	case ERR:
		c.Logger().Error("APPLOG: " + message)
	case FATAL:
		c.Logger().Error("APPLOG: " + message)
	default:
		c.Logger().Errorf("APPLOG, LEVEL_NOT_SET: %v", message)
	}

	log := &models.Log{
		SingleID: getCurrentSingle(c).ID,
		Category: cat,
		Level:    lev,
		Message:  message,
	}
	err := c.Value("tx").(*pop.Connection).Create(log)
	if err != nil {
		return err
	}
	return nil
}
