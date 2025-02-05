package stream

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DamirAbd/HLA-HW1/types"
	"github.com/DamirAbd/HLA-HW1/websockets"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	topic            string
	brokerAddress    string
	webSocketManager *websockets.Manager
}

func NewKafkaConsumer(topic, brokerAddress string, webSocketManager *websockets.Manager) *KafkaConsumer {
	return &KafkaConsumer{
		topic:            topic,
		brokerAddress:    brokerAddress,
		webSocketManager: webSocketManager,
	}
}

func (c *KafkaConsumer) Start() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.brokerAddress},
		Topic:   c.topic,
		//GroupID: "web-socket-group",
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message from Kafka: %v", err)
			continue
		}

		var post types.Post
		err = json.Unmarshal(msg.Value, &post)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		c.webSocketManager.BroadcastToAll(post)

	}
}
