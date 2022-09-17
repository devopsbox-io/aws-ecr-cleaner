package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	gerrors "github.com/pkg/errors"
)

func newAppRunnerClient() (*apprunner.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, gerrors.Wrapf(err, "cannot load aws config")
	}

	client := apprunner.NewFromConfig(cfg)

	return client, nil
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
