package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := "demo-topic"

	go produce(broker, topic)
	consume(broker, topic)
}

func produce(broker, topic string) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	defer w.Close()

	for i := 0; ; i++ {
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("key"),
				Value: []byte("Hello Kafka " + time.Now().String()),
			},
		)
		if err != nil {
			log.Println("produce error:", err)
		}

		time.Sleep(2 * time.Second)
	}
}

func consume(broker, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: "demo-group",
		Topic:   topic,
	})

	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("consume error:", err)
			continue
		}

		log.Printf("received: %s\n", string(msg.Value))
	}
}
