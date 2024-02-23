package tracker

import (
	"net/http"
)

type SettingsTemplateData struct {
	*CommonTemplateData
}

func (app *App) handleSettings(w http.ResponseWriter, r *http.Request) {
	tmpl := app.templator.Page("settings", nil)
	if err := tmpl.ExecuteTemplate(w, "layout", &SettingsTemplateData{
		CommonTemplateData: app.commonTemplateData(r),
	}); err != nil {
		app.internalError(w, err)
		return
	}
}
