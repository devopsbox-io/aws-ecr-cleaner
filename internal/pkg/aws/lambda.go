package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	gerrors "github.com/pkg/errors"
)

func newLambdaClient() (*lambda.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, gerrors.Wrapf(err, "cannot load aws config")
	}

	client := lambda.NewFromConfig(cfg)

	return client, nil
}

type LambdaClient interface {
	GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
}

type LambdaPaginators interface {
	NewListFunctionsPaginator(params *lambda.ListFunctionsInput, optFns ...func(*lambda.ListFunctionsPaginatorOptions)) LambdaListFunctionsPaginator
}

type lambdaPaginators struct {
	client *lambda.Client
}

func (l *lambdaPaginators) NewListFunctionsPaginator(params *lambda.ListFunctionsInput, optFns ...func(*lambda.ListFunctionsPaginatorOptions)) LambdaListFunctionsPaginator {
	return lambda.NewListFunctionsPaginator(l.client, params, optFns...)
}

type LambdaListFunctionsPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}
