package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

// Endless generator loop, say n = 500 emails, but each time it can only batch request up to 50 emails...

// 1. Create a channel that will store max of 50 emails at a time. We will do API calls and send result to this channel.
// 2. We should iterate over the channel with the 50 email objects, and do our processing, and then somehow trigger the next batch of 50 emails to be processed.
// 3. We should have a way to know when we are done processing all the emails, and then we can close the channel.

func main() {
	lazy_generator := make(chan string, 1)
	var abrupt_stop int32 = 0

	var batches_read int64 = 0
	go endless_mail_retriever(lazy_generator, &batches_read, &abrupt_stop)
	// We can have as many producers as we want - they will add the the generator channel. Note that the channel will only store max of 1 email batch at a time. This is to avoid preemptively storing all emails in memory.
	// go endless_mail_retriever(lazy_generator, &batches_read)

	for email_batch := range lazy_generator {
		if abrupt_stop == 1 {
			fmt.Println("Abrupt stop")
			break
		}

		process_emails(email_batch)
	}

	fmt.Printf("%d emails processed", batches_read)
}

// fetch 50 items -> process items -> fetch another 50 (or less) -> process items -> repeat until n total items processed.
func endless_mail_retriever(queue chan string, batches_read *int64, abrupt_stop *int32) {
	url := "https://jsonplaceholder.typicode.com/posts?_delay=1000&_limit=50"

	defer close(queue) // make sure to close the channel once all batches have been read

	for {
		if *batches_read == 3 {
			fmt.Println("Abrupt stop")
			atomic.StoreInt32(abrupt_stop, 1) // set abrupt stop to true (1)
			break
		}

		if *batches_read == 5 {
			fmt.Println("Done reading batches")
			break
		}

		// Make a GET request
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making the request:", err)
			return
		}
		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading the response body:", err)
			return
		}

		// Print the response body
		fmt.Println("Adding emails to queue...")

		atomic.AddInt64(batches_read, 1)

		queue <- string(body)
	}

}

func process_emails(emails string) {
	// fmt.Println("Processing emails:", emails)
}
