// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/aws/ssm.go

// Package aws is a generated GoMock package.
package aws

import (
	context "context"
	reflect "reflect"

	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	gomock "github.com/golang/mock/gomock"
)

// MockSsmPaginators is a mock of SsmPaginators interface.
type MockSsmPaginators struct {
	ctrl     *gomock.Controller
	recorder *MockSsmPaginatorsMockRecorder
}

// MockSsmPaginatorsMockRecorder is the mock recorder for MockSsmPaginators.
type MockSsmPaginatorsMockRecorder struct {
	mock *MockSsmPaginators
}

// NewMockSsmPaginators creates a new mock instance.
func NewMockSsmPaginators(ctrl *gomock.Controller) *MockSsmPaginators {
	mock := &MockSsmPaginators{ctrl: ctrl}
	mock.recorder = &MockSsmPaginatorsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSsmPaginators) EXPECT() *MockSsmPaginatorsMockRecorder {
	return m.recorder
}

// NewGetParametersByPathPaginator mocks base method.
func (m *MockSsmPaginators) NewGetParametersByPathPaginator(params *ssm.GetParametersByPathInput, optFns ...func(*ssm.GetParametersByPathPaginatorOptions)) SsmGetParametersByPathPaginator {
	m.ctrl.T.Helper()
	varargs := []interface{}{params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewGetParametersByPathPaginator", varargs...)
	ret0, _ := ret[0].(SsmGetParametersByPathPaginator)
	return ret0
}

// NewGetParametersByPathPaginator indicates an expected call of NewGetParametersByPathPaginator.
func (mr *MockSsmPaginatorsMockRecorder) NewGetParametersByPathPaginator(params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewGetParametersByPathPaginator", reflect.TypeOf((*MockSsmPaginators)(nil).NewGetParametersByPathPaginator), varargs...)
}

// MockSsmGetParametersByPathPaginator is a mock of SsmGetParametersByPathPaginator interface.
type MockSsmGetParametersByPathPaginator struct {
	ctrl     *gomock.Controller
	recorder *MockSsmGetParametersByPathPaginatorMockRecorder
}

// MockSsmGetParametersByPathPaginatorMockRecorder is the mock recorder for MockSsmGetParametersByPathPaginator.
type MockSsmGetParametersByPathPaginatorMockRecorder struct {
	mock *MockSsmGetParametersByPathPaginator
}

// NewMockSsmGetParametersByPathPaginator creates a new mock instance.
func NewMockSsmGetParametersByPathPaginator(ctrl *gomock.Controller) *MockSsmGetParametersByPathPaginator {
	mock := &MockSsmGetParametersByPathPaginator{ctrl: ctrl}
	mock.recorder = &MockSsmGetParametersByPathPaginatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSsmGetParametersByPathPaginator) EXPECT() *MockSsmGetParametersByPathPaginatorMockRecorder {
	return m.recorder
}

// HasMorePages mocks base method.
func (m *MockSsmGetParametersByPathPaginator) HasMorePages() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasMorePages")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasMorePages indicates an expected call of HasMorePages.
func (mr *MockSsmGetParametersByPathPaginatorMockRecorder) HasMorePages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasMorePages", reflect.TypeOf((*MockSsmGetParametersByPathPaginator)(nil).HasMorePages))
}

// NextPage mocks base method.
func (m *MockSsmGetParametersByPathPaginator) NextPage(ctx context.Context, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NextPage", varargs...)
	ret0, _ := ret[0].(*ssm.GetParametersByPathOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextPage indicates an expected call of NextPage.
func (mr *MockSsmGetParametersByPathPaginatorMockRecorder) NextPage(ctx interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextPage", reflect.TypeOf((*MockSsmGetParametersByPathPaginator)(nil).NextPage), varargs...)
}
