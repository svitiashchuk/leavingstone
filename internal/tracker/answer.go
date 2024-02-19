package tracker

import "net/http"

func (app *App) internalError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	app.errorLogger.Error(err.Error())
}

func (app *App) clientError(w http.ResponseWriter, err error) {
	http.Error(w, "Bad request", http.StatusBadRequest)
	app.appLogger.Error(err.Error())
}
