package main

import (
	"net/http"
	"ptocker/sqlite"
)

func main() {
	us, err := sqlite.NewUserService()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := us.Get("abigail.johnson@example.com")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(u.Name))
		}
	})

	err = http.ListenAndServe(":8887", nil)

	if err != nil {
		panic(err)
	}
}
