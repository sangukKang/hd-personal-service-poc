package downloader

import (
	"context"
	"encoding/json"

	"api-server-poc/components/messages"
	"api-server-poc/config"
	"api-server-poc/logger"
	"api-server-poc/utils"

	"github.com/segmentio/kafka-go"
)

type producerOutput struct {
	event *kafka.Writer
}

func NewProducerOutput(config config.KafkaWriteConfig) *producerOutput {
	writer, err := utils.NewKafkaProducer(&config)
	if err != nil {
		panic(err)
	}

	return &producerOutput{
		event: writer,
	}
}

// Download 요청 내용 kafka topic 에 전달
func (s *producerOutput) Stream(message messages.DownloadMessage) error {
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
