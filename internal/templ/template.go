package templ

import (
	"html/template"
	"leavingstone/internal/model"
	"math"
	"strings"
	"time"
	"unicode"
)

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
	}
}
