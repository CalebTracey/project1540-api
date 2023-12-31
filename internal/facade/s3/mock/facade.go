// Code generated by MockGen. DO NOT EDIT.
// Source: facade.go

// Package s3 is a generated GoMock package.
package s3

import (
	context "context"
	reflect "reflect"

	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	models "github.com/calebtracey/project1540-api/external/models"
	s30 "github.com/calebtracey/project1540-api/external/models/s3"
	gomock "go.uber.org/mock/gomock"
)

// MockIS3Facade is a mock of IS3Facade interface.
type MockIS3Facade struct {
	ctrl     *gomock.Controller
	recorder *MockIS3FacadeMockRecorder
}

// MockIS3FacadeMockRecorder is the mock recorder for MockIS3Facade.
type MockIS3FacadeMockRecorder struct {
	mock *MockIS3Facade
}

// NewMockIS3Facade creates a new mock instance.
func NewMockIS3Facade(ctrl *gomock.Controller) *MockIS3Facade {
	mock := &MockIS3Facade{ctrl: ctrl}
	mock.recorder = &MockIS3FacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIS3Facade) EXPECT() *MockIS3FacadeMockRecorder {
	return m.recorder
}

// DownloadS3Object mocks base method.
func (m *MockIS3Facade) DownloadS3Object(ctx context.Context, request s30.DownloadS3Request) (*s3.GetObjectOutput, *models.ErrorLog) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadS3Object", ctx, request)
	ret0, _ := ret[0].(*s3.GetObjectOutput)
	ret1, _ := ret[1].(*models.ErrorLog)
	return ret0, ret1
}

// DownloadS3Object indicates an expected call of DownloadS3Object.
func (mr *MockIS3FacadeMockRecorder) DownloadS3Object(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadS3Object", reflect.TypeOf((*MockIS3Facade)(nil).DownloadS3Object), ctx, request)
}

// GetS3ObjectNames mocks base method.
func (m *MockIS3Facade) GetS3ObjectNames(ctx context.Context, bucketName string) ([]string, *models.ErrorLog) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetS3ObjectNames", ctx, bucketName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(*models.ErrorLog)
	return ret0, ret1
}

// GetS3ObjectNames indicates an expected call of GetS3ObjectNames.
func (mr *MockIS3FacadeMockRecorder) GetS3ObjectNames(ctx, bucketName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetS3ObjectNames", reflect.TypeOf((*MockIS3Facade)(nil).GetS3ObjectNames), ctx, bucketName)
}

// UploadS3Object mocks base method.
func (m *MockIS3Facade) UploadS3Object(ctx context.Context, request s30.UploadS3Request) *models.ErrorLog {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadS3Object", ctx, request)
	ret0, _ := ret[0].(*models.ErrorLog)
	return ret0
}

// UploadS3Object indicates an expected call of UploadS3Object.
func (mr *MockIS3FacadeMockRecorder) UploadS3Object(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadS3Object", reflect.TypeOf((*MockIS3Facade)(nil).UploadS3Object), ctx, request)
}
