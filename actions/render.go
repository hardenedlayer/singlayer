package actions

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),

		// Add template helpers here:
		Helpers: render.Helpers{
			"timeYYMDHMS": func(t time.Time) string {
				return t.Local().Format("2006-01-02 15:04:05")
			},
			"timeYYMDHM": func(t time.Time) string {
				return t.Local().Format("2006-01-02 15:04")
			},
			"timeYMDHM": func(t time.Time) string {
				return t.Local().Format("06-01-02 15:04")
			},
			"timeMDHMS": func(t time.Time) string {
				return t.Local().Format("01-02 15:04:05")
			},
			"timeMDHM": func(t time.Time) string {
				return t.Local().Format("01-02 15:04")
			},
			"timeYYMD": func(t time.Time) string {
				return t.Local().Format("2006-01-02")
			},
			"timeYMD": func(t time.Time) string {
				return t.Local().Format("06-01-02")
			},
			"timeMD": func(t time.Time) string {
				return t.Local().Format("01-02")
			},
			"timeHMS": func(t time.Time) string {
				return t.Local().Format("15:04:05")
			},
			"timeHM": func(t time.Time) string {
				return t.Local().Format("15:04")
			},
			"iconize": func(s string) template.HTML {
				switch s {
				case "AUTO":
					return template.HTML(`<i class="fa fa-cog"></i>`)
				case "EMPLOYEE":
					return template.HTML(`<i class="fa fa-bars"></i>`)
				case "USER":
					return template.HTML(`<i class="fa fa-user-circle-o"></i>`)
				default:
					return template.HTML(`<i class="fa fa-` + s + `"></i>`)
				}
			},
			"shorten":  shortenHelper,
			"asArray":  asArrayHelper,
			"paginate": pagerHelper,
		},
	})
}

func shortenHelper(s string, l int) string {
	if len(s) > l {
		return s[0:l-3] + "..."
	} else {
		return s
	}
}

func asArrayHelper(a string) interface{} {
	sl := strings.Split(
		strings.Replace(
			strings.Trim(a, "[]"),
			" ", "", -1),
		",")
	switch len(sl) {
	case 0:
		return "none"
	case 1:
		return sl[0]
	default:
		label := fmt.Sprintf("%v and %v more", sl[0], len(sl)-1)
		return template.HTML(`<span data-toggle="tooltip" title="` +
			a + `">` + label + `</span>`)
	}
}

func pagerHelper(pos, pp, end int) template.HTML {
	var str string
	pager_len := 11
	center := pager_len/2 + 1
	arm := pager_len/2 - 2

	loop_start := 1
	loop_end := end
	fmt.Printf("pager: %v %v %v", pager_len, arm, center)

	if end > pager_len {
		loop_end = pager_len - 2
		if pos > center {
			loop_start = pos - arm
			loop_end = pos + arm
			str += fmt.Sprintf(`<li><a href="?page=1&pp=%v">1</a></li>`,
				pp)
			str += `<li><a>...</a></li>`
		}
		if pos > (end - arm - 3) {
			loop_end = end
			loop_start = end - pager_len + 3
		}
	}
	for i := loop_start; i <= loop_end; i++ {
		attr := ""
		if i == pos {
			attr = ` class="active"`
		}
		str += fmt.Sprintf(`<li%v><a href="?page=%v&pp=%v">%v</a></li>`,
			attr, i, pp, i)
	}
	if end > loop_end {
		str += `<li><a>...</a></li>`
		str += fmt.Sprintf(`<li><a href="?page=%v&pp=%v">%v</a></li>`,
			end, pp, end)
	}

	prev := pos - 1
	next := pos + 1
	prev_href := fmt.Sprintf(`?page=%v&pp=%v`, prev, pp)
	next_href := fmt.Sprintf(`?page=%v&pp=%v`, next, pp)
	prev_class := ""
	next_class := ""
	if next > end {
		next_class = "disabled"
		next_href = ""
	}
	if prev == 0 {
		prev_class = "disabled"
		prev_href = ""
	}

	return template.HTML(`<nav aria-label="Page navigation" class="text-center">
	<ul class="pagination">
		<li class="` + prev_class + `">
			<a href="` + prev_href + `" aria-label="Previous">
				<span aria-hidden="true">&laquo;</span>
			</a>
		</li>
` + str +
		`		<li class="` + next_class + `">
			<a href="` + next_href + `" aria-label="Next">
				<span aria-hidden="true">&raquo;</span>
			</a>
		</li>
	</ul>
</nav>`)
}
