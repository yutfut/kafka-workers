package kafka

import (
	"context"
	"fmt"

	"kafka-workers/conf"

	"github.com/segmentio/kafka-go"
)

func KafkaProducer(
	ctx context.Context,
	config *conf.Conf,
) (
	*kafka.Conn,
	error,
) {
	address := fmt.Sprintf("%s:%d", config.Kafka.Host, config.Kafka.Port)

	return kafka.DialLeader(
		ctx,
		"tcp",
		address,
		"input",
		0,
	)
}

func KafkaConsumer(
	config *conf.Conf,
) *kafka.Reader {
	address := fmt.Sprintf("%s:%d", config.Kafka.Host, config.Kafka.Port)

	consumer := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  []string{address},
			GroupID:  "consumer",
			Topic:    "output",
			MaxBytes: 10e6, // 10MB
		},
	)

	return consumer
}
