package main

// connect kafka using kgo

// create a record using kgo.Record{key:value}

// produceRecord:=client.ProduceSync(ctx,record)

// produceRecord.First()

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	seeds := []string{"kafka-service:9092"}

	var client *kgo.Client
	var err error
	//client, err = kgo.NewClient(    kgo.SeedBrokers(seeds...),    kgo.AllowAutoTopicCreation(),)

	for i := 0; i < 5; i++ {
		client, err = kgo.NewClient(
			kgo.SeedBrokers(seeds...), kgo.AllowAutoTopicCreation())
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
	record := &kgo.Record{Topic: "test-new-topic-1", Value: []byte("test message efef ")}

	ctx := context.Background()

	rec, err := client.ProduceSync(ctx, record).First()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Record produced successfully! Topic: %s, Partition: %d, Offset: %d, Timestamp: %s, Key: %v, Value: %s\n",
		rec.Topic, rec.Partition, rec.Offset, rec.Timestamp.Format(time.RFC3339), rec.Key, rec.Value)

}
