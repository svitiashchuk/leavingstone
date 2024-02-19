package handler

import "net/http"

func HTMXRedirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Add("HX-Redirect", url)
}
