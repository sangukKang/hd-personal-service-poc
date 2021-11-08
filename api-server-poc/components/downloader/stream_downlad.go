package downloader

import (
	"api-server-poc/components/messages"
)

type downloaderOutput struct {
}

func NewDownloaderOutput() *downloaderOutput {
	return &downloaderOutput{}
}

// S3 다운로드 스트림 처리
func (c *downloaderOutput) Stream(data messages.DownloadMessage) error {
	return nil
}
