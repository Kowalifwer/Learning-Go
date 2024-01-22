package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Endless generator loop, say n = 500 emails, but each time it can only batch request up to 50 emails...

// 1. Create a channel that will store max of 50 emails at a time. We will do API calls and send result to this channel.
// 2. We should iterate over the channel with the 50 email objects, and do our processing, and then somehow trigger the next batch of 50 emails to be processed.
// 3. We should have a way to know when we are done processing all the emails, and then we can close the channel.

func main() {
	lazy_generator := make(chan string, 1)
	var end_all int32 = 0
	var once sync.Once
	var wg sync.WaitGroup

	var batches_read int64 = 0

	N_PRODUCERS := 2

	for i := 0; i < N_PRODUCERS; i++ {
		wg.Add(1)
		go endless_mail_retriever(lazy_generator, &batches_read, &end_all, &once, &wg)
	}

	// go endless_mail_retriever(lazy_generator, &batches_read, &end_all, &once)
	// We can have as many producers as we want - they will add the the generator channel. Note that the channel will only store max of 1 email batch at a time. This is to avoid preemptively storing all emails in memory.
	// go endless_mail_retriever(lazy_generator, &batches_read)

	// time.Sleep(time.Second)
	go mail_processor(lazy_generator, &end_all)
	// go mail_processor(lazy_generator, &end_all)
	// go mail_processor(lazy_generator, &end_all)

	wg.Wait()

	fmt.Printf("%d emails processed", batches_read)
}

func mail_processor(queue chan string, end_all *int32) {
	fmt.Println("Starting consumer...")
	for email_batch := range queue {
		if atomic.LoadInt32(end_all) == 1 {
			fmt.Println("Abrupt stop")
			return
		}

		fmt.Println("CONSUMING from queue")
		process_emails(email_batch)
		time.Sleep(time.Millisecond * 200)
	}

	// for {
	// 	if *end_all == 1 {
	// 		fmt.Println("Abrupt stop")
	// 		return
	// 	}

	// 	select {
	// 	case email_batch, ok := <-queue:
	// 		if !ok {
	// 			fmt.Println("Channel closed")
	// 			return
	// 		}

	// 		// sleep 200 ms
	// 		time.Sleep(time.Millisecond * 200)
	// 		fmt.Println("CONSUMING from queue")
	// 		process_emails(email_batch)

	// 	default: //perhaps when queue is empty, but not closed.
	// 		// fmt.Println("No emails to process")

	// 	}
	// }
}

// fetch 50 items -> process items -> fetch another 50 (or less) -> process items -> repeat until n total items processed.
func endless_mail_retriever(queue chan string, batches_read *int64, end_all *int32, once *sync.Once, wg *sync.WaitGroup) {
	url := "https://jsonplaceholder.typicode.com/posts?_delay=1000&_limit=50"

	// make sure to close the channel once all batches have been read
	defer once.Do(func() {
		close(queue)
		atomic.StoreInt32(end_all, 1) // set abrupt stop to true (1)
	})

	defer wg.Done()

	for {
		// abrupt stop simulation
		// if *batches_read == 3 {
		// 	fmt.Println("Abrupt stop")
		// 	atomic.StoreInt32(end_all, 1) // set abrupt stop to true (1)
		// 	break
		// }

		// simulates end of emails
		if *batches_read == 20 {
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

		// send the email batch to the channel(blocking, if channel is full)
		if atomic.LoadInt32(end_all) == 1 {
			fmt.Println("all producers should stop - channel closed and not adding to queue")
			break
		}

		atomic.AddInt64(batches_read, 1)
		queue <- string(body)
		fmt.Println("PRODUCING to queue")
	}

}

func process_emails(emails string) {
	// TODO: process emails
}
