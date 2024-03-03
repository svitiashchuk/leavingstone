package tracker

import (
	"leavingstone/internal/session"
	"net/http"
)

func (app *App) handleThemeChange(w http.ResponseWriter, r *http.Request) {
	theme := r.PostFormValue("theme")
	app.sm.Get(r.Context().Value(session.SessionContextKey).(string)).Set("theme", theme)

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}
