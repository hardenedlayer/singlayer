package actions

import (
	"strings"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	lang := strings.ToLower(c.Value("lang").(string))
	return c.Render(200, r.HTML("index."+lang+".html"))
}
