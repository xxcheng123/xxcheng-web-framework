package xxcheng_web_framework

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FileUploader struct {
	FileField  string
	DstPatFunc func(file *multipart.FileHeader) string
}
type FileUploaderOption func(uploader *FileUploader)

func NewFileUploader(opts ...FileUploaderOption) *FileUploader {
	fu := &FileUploader{}
	for _, opt := range opts {
		opt(fu)
	}
	return fu
}

func (f *FileUploader) Handler() HandlerFunc {
	return func(ctx *Context) {
		file, fh, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败：读取错误")
			return
		}
		defer file.Close()
		dstFile, err := os.OpenFile(f.DstPatFunc(fh), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败：创建文件失败")
			return
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, file)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败：复制失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}
