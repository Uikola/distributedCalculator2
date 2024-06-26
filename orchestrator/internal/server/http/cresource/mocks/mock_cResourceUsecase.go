// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=mocks/mock_cResourceUsecase.go -package=mocks cResourceUseCase
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockcResourceUseCase is a mock of cResourceUseCase interface.
type MockcResourceUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockcResourceUseCaseMockRecorder
}

// MockcResourceUseCaseMockRecorder is the mock recorder for MockcResourceUseCase.
type MockcResourceUseCaseMockRecorder struct {
	mock *MockcResourceUseCase
}

// NewMockcResourceUseCase creates a new mock instance.
func NewMockcResourceUseCase(ctrl *gomock.Controller) *MockcResourceUseCase {
	mock := &MockcResourceUseCase{ctrl: ctrl}
	mock.recorder = &MockcResourceUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcResourceUseCase) EXPECT() *MockcResourceUseCaseMockRecorder {
	return m.recorder
}

// ListCResources mocks base method.
func (m *MockcResourceUseCase) ListCResources(ctx context.Context) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCResources", ctx)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCResources indicates an expected call of ListCResources.
func (mr *MockcResourceUseCaseMockRecorder) ListCResources(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCResources", reflect.TypeOf((*MockcResourceUseCase)(nil).ListCResources), ctx)
}
