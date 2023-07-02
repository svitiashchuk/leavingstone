package main

import (
	"fmt"
	httpp "ptocker/http"
)

func main() {
	s := httpp.NewServer()

	fmt.Println("Serving...")
	s.Serve()
}
