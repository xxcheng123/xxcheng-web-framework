package xxcheng_web_framework

import (
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestFileUploader_Handler(t *testing.T) {
	s := NewHTTPServer()

	s.POST("/upload", NewFileUploader(func(uploader *FileUploader) {
		uploader.FileField = "user_file"
	}, func(uploader *FileUploader) {
		uploader.DstPatFunc = func(fh *multipart.FileHeader) string {
			return filepath.Join("testdata", "uploads", fh.Filename)
		}
	}).Handler())

	s.Start(":9998")
}
