package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {

	seeds := []string{"kafka-service:9092"}

	var client *kgo.Client
	var err error

	client, err = kgo.NewClient(
		kgo.SeedBrokers(seeds...), kgo.ConsumerGroup("my-group-identifier"), kgo.ConsumeTopics("test-new-topic-1"), kgo.FetchMinBytes(1), kgo.FetchMaxWait(10*time.Millisecond))
	if err != nil {
		time.Sleep(2 * time.Second)
		return
	}

	defer client.Close()

	//kgo.ConsumerGroup("my-group-identifier"),kgo.ConsumeTopics(topic),kgo.FetchMinBytes(1),kgo.FetchMaxWait(10*time.Millisecond),

	ctx := context.Background()

	for {
		fetches := client.PollFetches(ctx)
		log.Println("fetched")

		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			time.Sleep(100 * time.Millisecond)
			fmt.Println(errs)
			continue
		}

		// We can iterate through a record iterator...
		iter := fetches.RecordIter()
		log.Println("got result in iter")
		for !iter.Done() {
			record := iter.Next()
			log.Println(string(record.Value), "from an iterator!")
		}

	}
}
