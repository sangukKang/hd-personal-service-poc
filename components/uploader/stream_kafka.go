package uploader

import (
	"context"
	"encoding/json"

	"api-server-poc/components/messages"
	"api-server-poc/config"
	"api-server-poc/logger"
	"api-server-poc/utils"

	"github.com/segmentio/kafka-go"
)

type kafkaOutput struct {
	event *kafka.Writer
}

func NewProducerOutput(config config.KafkaWriteConfig) *kafkaOutput {
	writer, err := utils.NewKafkaProducer(&config)
	if err != nil {
		panic(err)
	}

	return &kafkaOutput{
		event: writer,
	}
}

func (s *kafkaOutput) Stream(message messages.UploadMessage) error {
	jsonByte, _ := json.Marshal(message)
	logger.Info("send to kafka = ", string(jsonByte))
	err := s.event.WriteMessages(context.Background(), kafka.Message{
		Value: jsonByte,
	})
	if err != nil {
		logger.Error("err stream = ", err)
	}
	return err
}
