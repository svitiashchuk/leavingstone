package tracker

import (
	"html/template"
	"net/http"
)

type SettingsTemplateData struct {
	*CommonTemplateData
}

func (app *App) handleSettings(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"frontend/src/templates/layout.html",
		"frontend/src/templates/settings.html",
	))

	if err := tmpl.ExecuteTemplate(w, "layout", &SettingsTemplateData{
		CommonTemplateData: app.commonTemplateData(r),
	}); err != nil {
		app.internalError(w, err)
	}
}
