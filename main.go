// create function to print hello world
package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Set to run on port 3000, exposed via port forwarding
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
