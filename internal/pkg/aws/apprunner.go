package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
)

func newAppRunnerClient(cfg aws.Config) *apprunner.Client {
	return apprunner.NewFromConfig(cfg)
}

type AppRunnerClient interface {
	DescribeService(ctx context.Context, params *apprunner.DescribeServiceInput, optFns ...func(*apprunner.Options)) (*apprunner.DescribeServiceOutput, error)
}
type AppRunnerPaginators interface {
	NewListServicesPaginator(params *apprunner.ListServicesInput, optFns ...func(*apprunner.ListServicesPaginatorOptions)) AppRunnerListServicesPaginator
}

type appRunnerPaginators struct {
	client *apprunner.Client
}

func (a *appRunnerPaginators) NewListServicesPaginator(params *apprunner.ListServicesInput, optFns ...func(*apprunner.ListServicesPaginatorOptions)) AppRunnerListServicesPaginator {
	return apprunner.NewListServicesPaginator(a.client, params, optFns...)
}

type AppRunnerListServicesPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*apprunner.Options)) (*apprunner.ListServicesOutput, error)
}
