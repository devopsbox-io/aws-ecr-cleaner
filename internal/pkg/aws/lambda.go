package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func newLambdaClient(cfg aws.Config) *lambda.Client {
	return lambda.NewFromConfig(cfg)
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
