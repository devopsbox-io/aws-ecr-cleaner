package aws

import "github.com/golang/mock/gomock"

func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mockEcsClient := NewMockEcsClient(ctrl)
	mockEcsPaginators := NewMockEcsPaginators(ctrl)
	mockLambdaClient := NewMockLambdaClient(ctrl)
	mockLambdaPaginators := NewMockLambdaPaginators(ctrl)
	mockAppRunnerClient := NewMockAppRunnerClient(ctrl)
	mockAppRunnerPaginators := NewMockAppRunnerPaginators(ctrl)
	mockEcrClient := NewMockEcrClient(ctrl)
	mockEcrPaginators := NewMockEcrPaginators(ctrl)
	mockSsmPaginators := NewMockSsmPaginators(ctrl)

	return &MockProvider{
		Provider: &Provider{
			Region: "mock-aws-region",

			EcsClient:           mockEcsClient,
			EcsPaginators:       mockEcsPaginators,
			LambdaClient:        mockLambdaClient,
			LambdaPaginators:    mockLambdaPaginators,
			AppRunnerClient:     mockAppRunnerClient,
			AppRunnerPaginators: mockAppRunnerPaginators,
			EcrClient:           mockEcrClient,
			EcrPaginators:       mockEcrPaginators,
			SsmPaginators:       mockSsmPaginators,
		},
		MockEcsClient:           mockEcsClient,
		MockEcsPaginators:       mockEcsPaginators,
		MockLambdaClient:        mockLambdaClient,
		MockLambdaPaginators:    mockLambdaPaginators,
		MockAppRunnerClient:     mockAppRunnerClient,
		MockAppRunnerPaginators: mockAppRunnerPaginators,
		MockEcrClient:           mockEcrClient,
		MockEcrPaginators:       mockEcrPaginators,
		MockSsmPaginators:       mockSsmPaginators,
	}
}

type MockProvider struct {
	Provider *Provider

	MockEcsClient           *MockEcsClient
	MockEcsPaginators       *MockEcsPaginators
	MockLambdaClient        *MockLambdaClient
	MockLambdaPaginators    *MockLambdaPaginators
	MockAppRunnerClient     *MockAppRunnerClient
	MockAppRunnerPaginators *MockAppRunnerPaginators
	MockEcrClient           *MockEcrClient
	MockEcrPaginators       *MockEcrPaginators
	MockSsmPaginators       *MockSsmPaginators
}
