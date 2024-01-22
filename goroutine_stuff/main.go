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
	lazy_generator := make(chan string, 4)
	end_fetching := make(chan struct{})
	var once sync.Once
	var wg sync.WaitGroup

	var batches_read int64 = 0

	N_PRODUCERS := 10
	N_CONSUMERS := 1

	// Add to the wait group counter to ensure that the main at least waits for the PRODUCERS to finish.
	// Key idea of this solution, is our producers will increment the wait group counter ONCE they add work to the queue.
	// And our consumers will decrement the wait group counter ONCE they are done processing the work.
	// However, we need to have a preemptive wg.Add(1) to make sure the producers run. We also must make sure to wg.Done() after the producers are done (ONLY ONCE)
	wg.Add(1)

	for i := 0; i < N_PRODUCERS; i++ {
		go endless_mail_retriever(lazy_generator, &batches_read, end_fetching, &once, &wg, i+1)
	}

	// go endless_mail_retriever(lazy_generator, &batches_read, &end_fetching, &once)
	// We can have as many producers as we want - they will add the the generator channel. Note that the channel will only store max of 1 email batch at a time. This is to avoid preemptively storing all emails in memory.
	var batches_processed int64 = 0

	for i := 0; i < N_CONSUMERS; i++ {
		go mail_processor(lazy_generator, i+1, &batches_processed, &wg)
	}

	wg.Wait()

	fmt.Printf("%d email batches read and added to be proccessed \n", batches_read)
	fmt.Printf("%d email batches actually processed", batches_processed)
}

func mail_processor(queue chan string, id int, batches_processed *int64, wg *sync.WaitGroup) {
	fmt.Println("Starting consumer...")
	// RANGE LOOP OVER CHANNEL APPROACH (cleaner IMO)
	for email_batch := range queue {
		fmt.Println("CONSUMING from queue", id)
		process_emails(email_batch)
		atomic.AddInt64(batches_processed, 1)
		time.Sleep(time.Millisecond * 200)
		wg.Done()
	}

	// INFINITE LOOP WITH SELECT CASE APPROACH
	// for {
	// 	select {
	// 	case email_batch, ok := <-queue:
	// 		if !ok {
	// 			fmt.Println("Channel closed")
	// 			return
	// 		}

	// 		// sleep 200 ms
	// 		fmt.Println("CONSUMING from queue", id)
	// 		process_emails(email_batch)
	// 		time.Sleep(time.Millisecond * 200)
	// 		atomic.AddInt64(batches_processed, 1)
	// 		wg.Done()

	// 	default:
	// 		time.Sleep(time.Millisecond * 200) // sleep 200 ms if channel is empty but not closed (i.e wait for more work)
	// 	}
	// }
}

// fetch 50 items -> process items -> fetch another 50 (or less) -> process items -> repeat until n total items processed.
func endless_mail_retriever(queue chan string, batches_read *int64, end_fetching chan struct{}, once *sync.Once, wg *sync.WaitGroup, id int) {
	url := "https://jsonplaceholder.typicode.com/posts?_delay=1000&_limit=50"

	// make sure to close the channel once all batches have been read
	defer once.Do(func() {
		close(end_fetching)
		wg.Done()
	})

	for {
		// simulates end of emails
		if atomic.LoadInt64(batches_read) >= 6 {
			return
		}

		// Make a GET request
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making the request:", err)
			continue
		}
		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading the response body:", err)
			continue
		}

		select {
		case <-end_fetching:
			fmt.Println("Producer", id, "tried to add to queue after channel closed")
			return

		case queue <- string(body):
			fmt.Println("Producer", id, "added to queue")
			atomic.AddInt64(batches_read, 1)
			wg.Add(1)

		default:
			// fmt.Println("Channel full/attempt to add to queue failed")
		}

	}

}

func process_emails(emails string) {
	// TODO: process emails
}
