package uploader

import (
	"sync"

	"api-server-poc/components/messages"
	"api-server-poc/config"
)

type UploaderComponent struct {
	config  *config.Config
	pipline *Pipline
	once    sync.Once
}

func NewUploaderComponent(config *config.Config, pipline *Pipline) *UploaderComponent {
	return &UploaderComponent{
		config:  config,
		pipline: pipline,
	}
}

func (srv *UploaderComponent) Upload(data messages.UploadMessage) error {
	// 스트림 처리
	go srv.pipline.Output(data)
	return nil
}
