package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Include the navigation partial in the template files.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func longTask(ch chan<- string) {
	// Simulating a long task
	time.Sleep(5 * time.Second)
	ch <- "Long task completed"
}

func shortTask() string {
	// Simulating a short task
	return "Short task completed"
}

func longTaskHandler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan string)

	// Start the long task in a goroutine
	go longTask(ch)

	// Return a response immediately (Non BLOCKING) and the goroutine will continue the long task.
	fmt.Fprint(w, "Long task started. Check status later.")

	// You can optionally handle the result of the long task asynchronously
	go func() {
		result := <-ch
		// Process the result, update database, etc.
		fmt.Println(result)
	}()

	// BELOW IS BLOCKING WAITING FOR THE RESULT

	// // Wait for the long task to complete and get the result
	// result := <-ch

	// // Process the result, update database, etc.
	// fmt.Fprint(w, "\nResult: ", result)
}

func shortTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Perform the short task
	result := shortTask()

	// Return a response immediately
	fmt.Fprint(w, result)
}
