package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ptocker/internal/pkg/tracker"
	"ptocker/sqlite"
)

func main() {
	us, err := sqlite.NewUserService()
	if err != nil {
		panic(err)
	}

	ls, err := sqlite.NewLeaveService()
	if err != nil {
		panic(err)
	}

	t := tracker.NewTracker(us, ls)

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

	http.HandleFunc("/leaves/upcoming", func(w http.ResponseWriter, r *http.Request) {
		ull := t.UpcomingLeaves()

		list, err := json.Marshal(ull)
		if err != nil {
			panic(err)
		}

		w.Write(list)
	})

	fmt.Println("Listening on :8887")
	err = http.ListenAndServe(":8887", nil)

	if err != nil {
		panic(err)
	}
}
