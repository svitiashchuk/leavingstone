package main

import (
	"html/template"
	"net/http"
	"strconv"
)

/***********************
This is a simple demonstration of how to use the built-in template package in Go to implement
"template fragments" as described here: https://htmx.org/essays/template-fragments/
Go accomplishes this with the {{block}} action (described here: https://pkg.go.dev/text/template)
which defines and executes a template fragment inline inside of another template.  You only have
to wire up your application to use the correct template name and the fragment will be executed.
************************/

var page *template.Template

// init function sets up the template+fragment.  Most of the work is actually done here.
// In a larger program, this would likely be stored in a separate file, but this makes for a
// simple example.
func init() {

	page = template.New("main")

	page = template.Must(page.Parse(`<!DOCTYPE html>

	<html>
	<head>
		<script src="https://unpkg.com/htmx.org@1.8.0"></script>
		<link rel="stylesheet" href="https://the.missing.style"/>
		<title>Template Fragment Example</title>
	</head>
	<body>
		<h1>Template Fragment Example</h1>

		<p>This page demonstrates how to create and serve
		<a href="https://htmx.org/essays/template-fragments/">template fragments</a>
		using the <a href="https://pkg.go.dev/text/template">built-in template package</a> in Go.</p>

		<p>This is accomplished by using the "block" action in the template, which lets you
		define and execute a sub-template in a single step.</p>
		<!-- Here's the fragment.  We can target it by executing the "buttonOnly" template. -->
		{{block "buttonOnly" .}}
			<button hx-get="/?counter={{.next}}&template=buttonOnly" hx-swap="outerHTML">
				This Button Has Been Clicked {{.counter}} Times
			</button>
		{{end}}
	</body>
	</html>`))
}

// handleRequest does the work to execute the template (or fragment) and serve the result.
// It's mostly boilerplate, so don't get hung up on it.
func handleRequest(w http.ResponseWriter, r *http.Request) {

	// Collect state info to pass to the template
	counter, _ := strconv.Atoi(r.URL.Query().Get("counter"))
	templateName := r.URL.Query().Get("template")
	if templateName == "" {
		templateName = "main" // default value in case the query parameter is missing
	}

	// Pack state info into a map to pass to the template
	data := make(map[string]int)
	data["counter"] = counter
	data["next"] = counter + 1

	// Execute the template and handle errors
	if err := page.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// main is the entry point for the program. It sets up and executes the HTTP server.
func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}
