package kafka

import (
	"log/slog"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(topic string, brokers []string, log *slog.Logger) error {
	for _, broker := range brokers {
		log.Info("creating topic", slog.String("topic", topic), slog.String("broker", broker))
		conn, err := kafka.Dial("tcp", broker)
		if err != nil {
			log.Error("failed to connect to broker", slog.String("ERR", err.Error()))
			return err
		}
		log.Info("connected to broker", slog.String("broker", broker))
		defer conn.Close()

		controller, err := conn.Controller()
		conn.Close()
		if err != nil {
			log.Error("failed to get controller", slog.String("ERR", err.Error()))
			return err
		}
		log.Info("connected to controller", slog.String("controller", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))))

		ctrlConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			log.Error("failed to connect to controller", slog.String("ERR", err.Error()))
			return err
		}
		defer ctrlConn.Close()

		err = ctrlConn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err != nil {
			log.Error("failed to create topic", slog.String("ERR", err.Error()))
			return err
		}
		log.Info("topic created", slog.String("topic", topic), slog.String("broker", broker))
		ctrlConn.Close()
	}

	return nil
}
