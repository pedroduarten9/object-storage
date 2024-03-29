// Code generated by MockGen. DO NOT EDIT.
// Source: object-storage-gateway/internal/domain (interfaces: MinioLoadBalancer)

// Package domain is a generated GoMock package.
package domain

import (
	context "context"
	gateway "object-storage-gateway/internal/gateway"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMinioLoadBalancer is a mock of MinioLoadBalancer interface.
type MockMinioLoadBalancer struct {
	ctrl     *gomock.Controller
	recorder *MockMinioLoadBalancerMockRecorder
}

// MockMinioLoadBalancerMockRecorder is the mock recorder for MockMinioLoadBalancer.
type MockMinioLoadBalancerMockRecorder struct {
	mock *MockMinioLoadBalancer
}

// NewMockMinioLoadBalancer creates a new mock instance.
func NewMockMinioLoadBalancer(ctrl *gomock.Controller) *MockMinioLoadBalancer {
	mock := &MockMinioLoadBalancer{ctrl: ctrl}
	mock.recorder = &MockMinioLoadBalancerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMinioLoadBalancer) EXPECT() *MockMinioLoadBalancerMockRecorder {
	return m.recorder
}

// GetMinioClient mocks base method.
func (m *MockMinioLoadBalancer) GetMinioClient(arg0 context.Context, arg1, arg2 string) (gateway.Minio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMinioClient", arg0, arg1, arg2)
	ret0, _ := ret[0].(gateway.Minio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMinioClient indicates an expected call of GetMinioClient.
func (mr *MockMinioLoadBalancerMockRecorder) GetMinioClient(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinioClient", reflect.TypeOf((*MockMinioLoadBalancer)(nil).GetMinioClient), arg0, arg1, arg2)
}
