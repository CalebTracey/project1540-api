package facade

import "project1540-api/internal/facade/s3"

type IFacade interface {
	s3.IS3Facade
}

type Service struct {
	S3 s3.Service
}
