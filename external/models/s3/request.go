package s3

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type DownloadS3Request struct {
	FileName   string `json:"fileName,omitempty"`
	BucketName string `json:"bucketName,omitempty"`
}

func (r *DownloadS3Request) FromJSON(req *http.Request) DownloadS3Request {
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		log.Errorf("FromJSON: %v", err)
	}
	return *r
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

type RequestOptions func(*UploadS3Request)

func NewUploadS3Request(options RequestOptions) UploadS3Request {
	newRequest := new(UploadS3Request)
	options(newRequest)
	return *newRequest
}

func FromFile(header *multipart.FileHeader, info os.FileInfo, tempFile *os.File, dest string) RequestOptions {
	return func(r *UploadS3Request) {
		r.Name = header.Filename
		r.Size = strconv.Itoa(int(info.Size()))
		r.DestBucket = dest
		r.Type = header.Header.Get("Content-Type")
		r.File = tempFile
	}
}

type Tags []Tag

type Tag struct {
	Name string `json:"name"`
}
