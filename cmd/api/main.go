package main

import (
	"fmt"
	"net/http"
	"ptocker/sqlite"
)

func main() {
	us, err := sqlite.NewUserService()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := us.Find("abigail.johnson@example.com")
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(u.Name))
		}
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		uu, err := us.AllUsers()
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(fmt.Sprint(uu)))
		}

		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Listening on :8887")
	err = http.ListenAndServe(":8887", nil)

	if err != nil {
		panic(err)
	}
}
