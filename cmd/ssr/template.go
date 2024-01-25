package main

import (
	"html/template"
	"leavingstone"
	"math"
	"strings"
	"time"
	"unicode"
)

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"humanDate": func(t time.Time) string {
			return t.Format("02 Jan")
		},
		"leaveDays": func(l leavingstone.Leave) int {
			return int(math.Round(l.Duration().Hours() / 24))
		},
		"leaveTypeSign": func(l leavingstone.Leave) string {
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
		"leaveTypeColor": func(l leavingstone.Leave) string {
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
		"monthNum": func(m time.Month) int {
			return int(m)
		},
	}
}
