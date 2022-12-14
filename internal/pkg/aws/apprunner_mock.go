// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/aws/apprunner.go

// Package aws is a generated GoMock package.
package aws

import (
	context "context"
	reflect "reflect"

	apprunner "github.com/aws/aws-sdk-go-v2/service/apprunner"
	gomock "github.com/golang/mock/gomock"
)

// MockAppRunnerClient is a mock of AppRunnerClient interface.
type MockAppRunnerClient struct {
	ctrl     *gomock.Controller
	recorder *MockAppRunnerClientMockRecorder
}

// MockAppRunnerClientMockRecorder is the mock recorder for MockAppRunnerClient.
type MockAppRunnerClientMockRecorder struct {
	mock *MockAppRunnerClient
}

// NewMockAppRunnerClient creates a new mock instance.
func NewMockAppRunnerClient(ctrl *gomock.Controller) *MockAppRunnerClient {
	mock := &MockAppRunnerClient{ctrl: ctrl}
	mock.recorder = &MockAppRunnerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppRunnerClient) EXPECT() *MockAppRunnerClientMockRecorder {
	return m.recorder
}

// DescribeService mocks base method.
func (m *MockAppRunnerClient) DescribeService(ctx context.Context, params *apprunner.DescribeServiceInput, optFns ...func(*apprunner.Options)) (*apprunner.DescribeServiceOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeService", varargs...)
	ret0, _ := ret[0].(*apprunner.DescribeServiceOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeService indicates an expected call of DescribeService.
func (mr *MockAppRunnerClientMockRecorder) DescribeService(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeService", reflect.TypeOf((*MockAppRunnerClient)(nil).DescribeService), varargs...)
}

// MockAppRunnerPaginators is a mock of AppRunnerPaginators interface.
type MockAppRunnerPaginators struct {
	ctrl     *gomock.Controller
	recorder *MockAppRunnerPaginatorsMockRecorder
}

// MockAppRunnerPaginatorsMockRecorder is the mock recorder for MockAppRunnerPaginators.
type MockAppRunnerPaginatorsMockRecorder struct {
	mock *MockAppRunnerPaginators
}

// NewMockAppRunnerPaginators creates a new mock instance.
func NewMockAppRunnerPaginators(ctrl *gomock.Controller) *MockAppRunnerPaginators {
	mock := &MockAppRunnerPaginators{ctrl: ctrl}
	mock.recorder = &MockAppRunnerPaginatorsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppRunnerPaginators) EXPECT() *MockAppRunnerPaginatorsMockRecorder {
	return m.recorder
}

// NewListServicesPaginator mocks base method.
func (m *MockAppRunnerPaginators) NewListServicesPaginator(params *apprunner.ListServicesInput, optFns ...func(*apprunner.ListServicesPaginatorOptions)) AppRunnerListServicesPaginator {
	m.ctrl.T.Helper()
	varargs := []interface{}{params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewListServicesPaginator", varargs...)
	ret0, _ := ret[0].(AppRunnerListServicesPaginator)
	return ret0
}

// NewListServicesPaginator indicates an expected call of NewListServicesPaginator.
func (mr *MockAppRunnerPaginatorsMockRecorder) NewListServicesPaginator(params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewListServicesPaginator", reflect.TypeOf((*MockAppRunnerPaginators)(nil).NewListServicesPaginator), varargs...)
}

// MockAppRunnerListServicesPaginator is a mock of AppRunnerListServicesPaginator interface.
type MockAppRunnerListServicesPaginator struct {
	ctrl     *gomock.Controller
	recorder *MockAppRunnerListServicesPaginatorMockRecorder
}

// MockAppRunnerListServicesPaginatorMockRecorder is the mock recorder for MockAppRunnerListServicesPaginator.
type MockAppRunnerListServicesPaginatorMockRecorder struct {
	mock *MockAppRunnerListServicesPaginator
}

// NewMockAppRunnerListServicesPaginator creates a new mock instance.
func NewMockAppRunnerListServicesPaginator(ctrl *gomock.Controller) *MockAppRunnerListServicesPaginator {
	mock := &MockAppRunnerListServicesPaginator{ctrl: ctrl}
	mock.recorder = &MockAppRunnerListServicesPaginatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppRunnerListServicesPaginator) EXPECT() *MockAppRunnerListServicesPaginatorMockRecorder {
	return m.recorder
}

// HasMorePages mocks base method.
func (m *MockAppRunnerListServicesPaginator) HasMorePages() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasMorePages")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasMorePages indicates an expected call of HasMorePages.
func (mr *MockAppRunnerListServicesPaginatorMockRecorder) HasMorePages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasMorePages", reflect.TypeOf((*MockAppRunnerListServicesPaginator)(nil).HasMorePages))
}

// NextPage mocks base method.
func (m *MockAppRunnerListServicesPaginator) NextPage(ctx context.Context, optFns ...func(*apprunner.Options)) (*apprunner.ListServicesOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NextPage", varargs...)
	ret0, _ := ret[0].(*apprunner.ListServicesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextPage indicates an expected call of NextPage.
func (mr *MockAppRunnerListServicesPaginatorMockRecorder) NextPage(ctx interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextPage", reflect.TypeOf((*MockAppRunnerListServicesPaginator)(nil).NextPage), varargs...)
}
