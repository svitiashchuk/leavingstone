package main

import (
	"fmt"
	httpp "leavingstone/http"
)

func main() {
	s := httpp.NewServer()

	addr := ":8080"
	fmt.Printf("Serving... at %s\n", addr)
	s.Serve(addr)
}
