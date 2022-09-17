package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	gerrors "github.com/pkg/errors"
)

func newEcrClient() (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, gerrors.Wrapf(err, "cannot load aws config")
	}

	client := ecr.NewFromConfig(cfg)

	return client, nil
}

type EcrClient interface {
	ListTagsForResource(ctx context.Context, params *ecr.ListTagsForResourceInput, optFns ...func(*ecr.Options)) (*ecr.ListTagsForResourceOutput, error)
	BatchDeleteImage(ctx context.Context, params *ecr.BatchDeleteImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchDeleteImageOutput, error)
}

type EcrPaginators interface {
	NewDescribeRepositoriesPaginator(params *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.DescribeRepositoriesPaginatorOptions)) EcrDescribeRepositoriesPaginator
	NewDescribeImagesPaginator(params *ecr.DescribeImagesInput, optFns ...func(*ecr.DescribeImagesPaginatorOptions)) EcrDescribeImagesPaginator
}

type ecrPaginators struct {
	client *ecr.Client
}

func (e *ecrPaginators) NewDescribeRepositoriesPaginator(params *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.DescribeRepositoriesPaginatorOptions)) EcrDescribeRepositoriesPaginator {
	return ecr.NewDescribeRepositoriesPaginator(e.client, params, optFns...)
}

type EcrDescribeRepositoriesPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
}

func (e *ecrPaginators) NewDescribeImagesPaginator(params *ecr.DescribeImagesInput, optFns ...func(*ecr.DescribeImagesPaginatorOptions)) EcrDescribeImagesPaginator {
	return ecr.NewDescribeImagesPaginator(e.client, params, optFns...)
}

type EcrDescribeImagesPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
}
