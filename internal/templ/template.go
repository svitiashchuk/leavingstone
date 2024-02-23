package templ

import (
	"html/template"
	"leavingstone/internal/model"
	"math"
	"strings"
	"time"
	"unicode"
)

type Templator struct {
	pagesCache     map[string]*template.Template
	fragmentsCache map[string]*template.Template
}

func New() *Templator {
	return &Templator{
		pagesCache:     make(map[string]*template.Template),
		fragmentsCache: make(map[string]*template.Template),
	}
}

func (t *Templator) Fragment(name string, tFuncs template.FuncMap) *template.Template {
	if _, ok := t.fragmentsCache[name]; !ok {
		filename := "frontend/src/templates/fragments/" + name + ".html"

		tem := template.New(name)
		if len(tFuncs) > 0 {
			tem.Funcs(tFuncs)
		}
		t.fragmentsCache[name] = template.Must(tem.ParseFiles(filename))
	}

	return t.fragmentsCache[name]
}

func (t *Templator) Page(name string, tFuncs template.FuncMap) *template.Template {
	if _, ok := t.pagesCache[name]; !ok {
		pageHTMLFilename := "frontend/src/templates/pages/" + name + ".html"
		files := append(basicTemplates(), pageHTMLFilename)

		tem := template.New(name)
		if len(tFuncs) > 0 {
			tem.Funcs(tFuncs)
		}

		t.pagesCache[name] = template.Must(
			tem.ParseFiles(files...),
		)
	}

	return t.pagesCache[name]
}

func basicTemplates() []string {
	return []string{
		"frontend/src/templates/layout.html",
		"frontend/src/templates/partials/nav.html",
		"frontend/src/templates/partials/alert.html",
		"frontend/src/templates/partials/sidebar.html",
	}
}

func Funcs() template.FuncMap {
	return template.FuncMap{
		"humanDate": func(t time.Time) string {
			return t.Format("02 Jan")
		},
		"leaveTypeSign": func(l model.Leave) string {
			if l.Type == "vacation" {
				return "✈"
			}
			if l.Type == "sick" {
				return "✚"
			}
			if l.Type == "dayoff" {
				return "⧗"
			}

			panic("unknown leave type")
		},
		"leaveTypeColor": func(l model.Leave) string {
			if l.Type == "vacation" {
				return "accent"
			}
			if l.Type == "sick" {
				return "warning"
			}
			if l.Type == "dayoff" {
				return "primary"
			}

			panic("unknown leave type")
		},
		"calculatePercentage": func(x, total int) int {
			return int(math.Round(float64(x) / float64(total) * 100))
		},
		"nameAbbrev": func(s string) string {
			return strings.Map(func(r rune) rune {
				if unicode.IsUpper(r) {
					return r
				} else {
					return -1
				}
			}, s)
		},
		"firstName": func(s string) string {
			return strings.Split(s, " ")[0]
		},
		"lastName": func(s string) string {
			return strings.Join(strings.Split(s, " ")[1:], " ")
		},
		"monthNum": func(m time.Month) int {
			return int(m)
		},
		"availabilityBadge": func(leaveType string) string {
			if leaveType == "vacation" {
				return `<div class="badge badge-accent">Vacation</div>`
			}
			if leaveType == "sick" {
				return `<div class="badge badge-warning">Sick Leave</div>`
			}
			if leaveType == "dayoff" {
				return `<div class="badge badge-primary">Day Off</div>`
			}

			// TODO handle weekend / bank holiday
			return `<div class="badge">Available</div>`
		},
	}
}
