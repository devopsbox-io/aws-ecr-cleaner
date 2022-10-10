package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func newSsmClient(cfg aws.Config) *ssm.Client {
	return ssm.NewFromConfig(cfg)
}

type SsmPaginators interface {
	NewGetParametersByPathPaginator(params *ssm.GetParametersByPathInput, optFns ...func(*ssm.GetParametersByPathPaginatorOptions)) SsmGetParametersByPathPaginator
}

type ssmPaginators struct {
	client *ssm.Client
}

func (s *ssmPaginators) NewGetParametersByPathPaginator(params *ssm.GetParametersByPathInput, optFns ...func(*ssm.GetParametersByPathPaginatorOptions)) SsmGetParametersByPathPaginator {
	return ssm.NewGetParametersByPathPaginator(s.client, params, optFns...)
}

type SsmGetParametersByPathPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
}
