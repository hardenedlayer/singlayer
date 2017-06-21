package actions

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
)

func AuthenticateHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user_id := c.Session().Get("user_id")
		if user_id == nil {
			c.Flash().Add("danger", "Sorry, Login required.")
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		err := next(c)
		return err
	}
}

func AdminPageKeeper(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		is_admin := c.Session().Get("is_admin")
		if is_admin == false {
			c.Flash().Add("danger", "STAFF ONLY")
			l(c, VIOL, FATAL, "user tried to enter god place!")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		c.Set("theme", "admin")
		err := next(c)
		return err
	}
}

func LanguageHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// quick and dirty, static ordered list of supported languages.
		supported_langs := []string{"ko-KR", "en-US"}
		langs := c.Request().Header.Get("Accept-Language")
		var lang string
		for _, la := range strings.Split(langs, ",") {
			la = strings.Split(la, ";")[0]
			for _, al := range supported_langs {
				if la == al {
					lang = la
					break
				}
			}
			if len(lang) > 0 {
				break
			}
		}
		if len(lang) == 0 {
			lang = supported_langs[0]
		}
		c.Set("lang", lang)
		c.Logger().Infof("current language: %v", c.Value("lang"))
		err := next(c)
		return err
	}
}

// check permission for specific pathes.
func PermissionHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		is_admin := c.Session().Get("is_admin")
		if is_admin == nil {
			c.Logger().Errorf("cannot get session info for PermissionHandler")
			return c.Error(500, errors.New("Session Information Error"))
		}
		if is_admin.(bool) == false {
			pos := strings.Split(c.Value("current_path").(string), "/")[1]
			perms := c.Session().Get("permissions").(string)

			// register pages requiring specific permission:
			perm := map[string]string{
				"landscape":   "landscape",
				"computes":    "landscape",
				"tickets":     "ticket",
				"directlinks": "dlink",
				"docs":        "user",
			}
			if p := perm[pos]; p != "" {
				if strings.Contains(perms, p) == false {
					l(c, VIOL, FATAL, "violation no perm_%v for %v", p, pos)
					c.Flash().Add("danger",
						"You don't have permission for "+pos+"!")
					return c.Redirect(http.StatusTemporaryRedirect, "/")
				}
				c.Logger().Debugf("user aquires permission %v for %v", p, pos)
			}
		}
		err := next(c)
		return err
	}
}

// call by every page requests
func SessionInfoHandler(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user_id := c.Session().Get("user_id")
		if user_id != nil {
			c.Set("user_id", user_id)
			c.Set("user_name", c.Session().Get("user_name"))
			c.Set("user_mail", c.Session().Get("user_mail"))
			c.Set("user_icon", c.Session().Get("user_icon"))
			c.Set("user_is_admin", c.Session().Get("is_admin"))
			c.Set("actors", c.Session().Get("actors"))
		}
		if perms, ok := c.Session().Get("permissions").(string); ok {
			for _, el := range strings.Split(perms, ":") {
				if el != "" {
					c.Set("perm_"+el, true)
				}
			}
		}
		c.Set("actor", "All")
		if k, err := c.Request().Cookie("_singlayer_actor"); err == nil {
			if actors, ok := c.Session().Get("actors").([]string); ok {
				for _, v := range actors {
					if v == k.Value {
						c.Set("actor", k.Value)
					}
				}
			}
		}
		err := next(c)
		return err
	}
}
