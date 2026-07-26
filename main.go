package main

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	http.ListenAndServe(":8080", nil)
}
