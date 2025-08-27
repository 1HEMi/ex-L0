package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	brokers := flag.String("brokers", "127.0.0.1:9092", "comma-separated brokers")
	topic := flag.String("topic", "orders", "topic")
	file := flag.String("file", "scripts/seed_order.json", "json file")
	flag.Parse()

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  strings.Split(*brokers, ","),
		Topic:    *topic,
		Balancer: &kafka.Hash{}})
	defer w.Close()

	data, err := os.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	if err := w.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
		log.Fatal(err)
	}
	log.Println("message produced")
}
