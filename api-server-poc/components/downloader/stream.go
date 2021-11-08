package downloader

import (
	"reflect"
	"sync"

	"api-server-poc/components/messages"
	"api-server-poc/config"
	"api-server-poc/logger"
)

type Pipline struct {
	config  *config.Config
	streams []StreamOutput
	once    sync.Once
}

func NewPipLine(config *config.Config) *Pipline {
	pipline := &Pipline{
		config: config,
	}
	pipline.setStream(
		NewDownloaderOutput(),
		NewProducerOutput(config.KafkaConfig.ProducerEvent),
	)
	return pipline
}

func (p *Pipline) setStream(streams ...StreamOutput) {
	for _, stream := range streams {
		p.streams = append(p.streams, stream)
	}
}

func (p *Pipline) Output(data messages.DownloadMessage) {
	logger.Debug(data)
	for _, output := range p.streams {
		go func(o StreamOutput) {
			err := o.Stream(data)
			if err != nil {
				logger.Errorf("err stream = %v, reason = %v ", reflect.TypeOf(output), err)
			}
		}(output)
	}
}

type StreamOutput interface {
	Stream(message messages.DownloadMessage) error
}
