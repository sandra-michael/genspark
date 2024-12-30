package kafka

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Conf struct {
	client *kgo.Client
	admin  *kadm.Client
}

func NewConf() (*Conf, error) {
	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")

	if host == "" || port == "" {
		return nil, fmt.Errorf("kafka host or port is empty")
	}
	connString := fmt.Sprintf("%s:%s", host, port)
	var err error
	var client *kgo.Client
	for i := 1; i < 8; i++ {

		var backoff time.Duration = 2
		client, err = kgo.NewClient(
			kgo.SeedBrokers(connString),

			//ProducerLinger sets how long individual topic partitions will linger waiting for more records
			//before triggering a request to be built.
			kgo.ProducerLinger(0),
			kgo.AllowAutoTopicCreation(),
		)
		if err != nil {
			fmt.Printf("kafka client error: %v", err)
			time.Sleep(backoff * time.Second)
			backoff *= 2
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		err = client.Ping(ctx)

		if err != nil {
			fmt.Printf("kafka client error: %v", err)
			time.Sleep(backoff * time.Second)
			backoff *= 2
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("kafka client error: %v", err)
	}
	admin := kadm.NewClient(client)
	return &Conf{
		client: client,
		admin:  admin,
	}, nil
}

func (c *Conf) ProduceMessage(ctx context.Context, topicName string, key []byte, value []byte) error {

	record := &kgo.Record{Topic: topicName, Key: key, Value: value}

	rec, err := c.client.ProduceSync(ctx, record).First()

	if err != nil {
		return fmt.Errorf("kafka producer error: %v", err)
	}

	fmt.Printf("Record produced successfully! Topic: %s, Partition: %d, Offset: %d, Timestamp: %s, Key: %v, Value: %s\n",
		rec.Topic, rec.Partition, rec.Offset, rec.Timestamp.Format(time.RFC3339), rec.Key, rec.Value)

	return nil

}
