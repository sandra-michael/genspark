// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source service.go -destination mockmodels/service_mock.go -package mockmodels
//

// Package mockmodels is a generated GoMock package.
package mockmodels

import (
	models "Assignment1/models"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateTable mocks base method.
func (m *MockService) CreateTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockServiceMockRecorder) CreateTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockService)(nil).CreateTable))
}

// CreateTask mocks base method.
func (m *MockService) CreateTask(ctx context.Context, newTask models.NewTask) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, newTask)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockServiceMockRecorder) CreateTask(ctx, newTask any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockService)(nil).CreateTask), ctx, newTask)
}

// DeleteTask mocks base method.
func (m *MockService) DeleteTask(ctx context.Context, taskID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockServiceMockRecorder) DeleteTask(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockService)(nil).DeleteTask), ctx, taskID)
}

// FetchTask mocks base method.
func (m *MockService) FetchTask(ctx context.Context, id int) (models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchTask", ctx, id)
	ret0, _ := ret[0].(models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTask indicates an expected call of FetchTask.
func (mr *MockServiceMockRecorder) FetchTask(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTask", reflect.TypeOf((*MockService)(nil).FetchTask), ctx, id)
}

// FetchTasks mocks base method.
func (m *MockService) FetchTasks(ctx context.Context) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchTasks", ctx)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTasks indicates an expected call of FetchTasks.
func (mr *MockServiceMockRecorder) FetchTasks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTasks", reflect.TypeOf((*MockService)(nil).FetchTasks), ctx)
}

// Ping mocks base method.
func (m *MockService) Ping() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Ping")
}

// Ping indicates an expected call of Ping.
func (mr *MockServiceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockService)(nil).Ping))
}

// UpdateTask mocks base method.
func (m *MockService) UpdateTask(ctx context.Context, id int, updateTask models.UpdateTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, id, updateTask)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockServiceMockRecorder) UpdateTask(ctx, id, updateTask any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockService)(nil).UpdateTask), ctx, id, updateTask)
}

// UpdateTaskStatus mocks base method.
func (m *MockService) UpdateTaskStatus(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskStatus", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskStatus indicates an expected call of UpdateTaskStatus.
func (mr *MockServiceMockRecorder) UpdateTaskStatus(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskStatus", reflect.TypeOf((*MockService)(nil).UpdateTaskStatus), ctx, id)
}
