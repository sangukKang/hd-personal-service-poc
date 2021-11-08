package uploader

import (
	"api-server-poc/components/messages"
)

type uploaderOutput struct {
}

func NewUploaderOutput() *uploaderOutput {

	return &uploaderOutput{}
}

// S3 에 업로드 스트림 처리
func (c *uploaderOutput) Stream(data messages.UploadMessage) error {

	return nil
}
