package main

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

func main() {
	// Set up a new Kafka client
	seeds := []string{"kafka-service:9092"}
	var client *kgo.Client
	var err error
	for i := 0; i < 5; i++ {
		client, err = kgo.NewClient(
			kgo.SeedBrokers(seeds...))
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		err = client.Ping(context.Background())
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
	}
	if err != nil {
		panic(err)
	}
	defer client.Close()
	adminClient := kadm.NewClient(client)
	err = CreatTopic(adminClient, "test-new-topic-1")
	if err != nil {
		panic(err)
	}

}

func CreatTopic(adminClient *kadm.Client, topic string) (err error) {
	// The number of partitions for the topic.
	// Partitions allow parallelism in Kafka. Minimum is 1.
	var partitions int32 = 1

	// The replication factor for the topic.
	// This is the number of copies of the data to ensure fault tolerance. Set it to at least 1.
	var replicationFactor int16 = 1
	ctx := context.Background()
	results, err := adminClient.CreateTopic(ctx, partitions, replicationFactor, nil, topic)
	if err != nil {
		return err
	}
	fmt.Println(results.Topic, "created")
	return nil

}
