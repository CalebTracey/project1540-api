package s3

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

type DownloadS3Request struct {
	FileName   string `json:"fileName,omitempty"`
	BucketName string `json:"bucketName,omitempty"`
}

type UploadS3Request struct {
	ID         string        `json:"id"`
	Name       string        `json:"name,omitempty"`
	DestBucket string        `json:"destBucket"`
	Size       string        `json:"size,omitempty"`
	Type       string        `json:"type,omitempty"`
	URL        string        `json:"url,omitempty"`
	Tags       Tags          `json:"tags,omitempty"`
	File       io.ReadCloser `json:"file,omitempty"`
}

type RequestMapper func(*UploadS3Request)

func NewUploadS3Request(req RequestMapper) UploadS3Request {
	newRequest := new(UploadS3Request)

	req(newRequest)

	return *newRequest
}

func FromFile(header *multipart.FileHeader, info os.FileInfo, tempFile *os.File, dest string) RequestMapper {
	return func(r *UploadS3Request) {
		r.Name = header.Filename
		r.Size = strconv.Itoa(int(info.Size()))
		r.DestBucket = dest
		r.Type = header.Header.Get("Content-Type")
		r.File = tempFile
	}
}

//func (r *UploadS3Request) FromFile(header *multipart.FileHeader, info os.FileInfo, tempFile *os.File, dest string) *UploadS3Request {
//r.Name = header.Filename
//r.
//	return r
//}

type Tags []Tag

type Tag struct {
	Name string `json:"name"`
}
