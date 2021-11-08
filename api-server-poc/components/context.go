package components

import (
	"api-server-poc/components/downloader"
	"api-server-poc/components/uploader"
	"api-server-poc/config"
)

type UploadServiceContext struct {
	config   *config.Config
	uploader *uploader.UploaderComponent
}

func NewUploadServiceContext(config *config.Config, uploader *uploader.UploaderComponent) *UploadServiceContext {
	return &UploadServiceContext{
		config:   config,
		uploader: uploader,
	}
}

type DownloadServiceContext struct {
	config     *config.Config
	downloader *downloader.DownloaderComponent
}

func NewDownloadServiceContext(config *config.Config, downloader *downloader.DownloaderComponent) *DownloadServiceContext {
	return &DownloadServiceContext{
		config:     config,
		downloader: downloader,
	}
}
