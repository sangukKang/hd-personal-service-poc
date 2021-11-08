package downloader

import (
	"sync"

	"api-server-poc/components/messages"
	"api-server-poc/config"
)

type DownloaderComponent struct {
	config  *config.Config
	pipline *Pipline
	once    sync.Once
}

func NewDownloaderComponent(config *config.Config, pipline *Pipline) *DownloaderComponent {
	return &DownloaderComponent{
		config:  config,
		pipline: pipline,
	}
}

func (srv *DownloaderComponent) Download(data messages.DownloadMessage) error {
	// 스트림 처리
	go srv.pipline.Output(data)
	return nil
}
