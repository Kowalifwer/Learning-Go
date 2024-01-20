// create function to print hello world
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Hello World</h1>")
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>About</h1>")
	})

	fmt.Println("Server starting...")
	http.ListenAndServe(":3000", nil)
}

// run the program
