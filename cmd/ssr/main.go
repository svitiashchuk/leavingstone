package main

import (
	"fmt"
	httpp "leavingstone/http"
)

func main() {
	s := httpp.NewServer()

	fmt.Println("Serving...")
	s.Serve()
}
