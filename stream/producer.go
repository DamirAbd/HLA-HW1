package stream

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/DamirAbd/HLA-HW1/types"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokerAddress, topic string) (*KafkaProducer, error) {
	if brokerAddress == "" || topic == "" {
		return nil, errors.New("broker address and topic must not be empty")
	}
	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokerAddress),
		Topic: topic,
	}

	return &KafkaProducer{writer: writer}, nil
}

func (p *KafkaProducer) SendPost(post types.Post) error {
	payload, err := json.Marshal(post)
	if err != nil {
		log.Printf("Failed to marshal post: %v", err)
		return err
	}

	err = p.writer.WriteMessages(context.Background(), kafka.Message{
		Value: payload,
	})
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}

	log.Println("Message sent to Kafka successfully!")
	return nil
}

func (p *KafkaProducer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Failed to close Kafka writer: %v", err)
	}
}
