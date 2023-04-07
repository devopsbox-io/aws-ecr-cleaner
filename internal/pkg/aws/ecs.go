package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func newEcsClient(cfg aws.Config) *ecs.Client {
	return ecs.NewFromConfig(cfg)
}

type EcsClient interface {
	DescribeServices(ctx context.Context, params *ecs.DescribeServicesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
	DescribeTaskDefinition(ctx context.Context, params *ecs.DescribeTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
}

type EcsPaginators interface {
	NewListClustersPaginator(params *ecs.ListClustersInput, optFns ...func(*ecs.ListClustersPaginatorOptions)) EcsListClustersPaginator
	NewListServicesPaginator(params *ecs.ListServicesInput, optFns ...func(*ecs.ListServicesPaginatorOptions)) EcsListServicesPaginator
	NewListTaskDefinitionsPaginator(param *ecs.ListTaskDefinitionsInput, optFns ...func(options *ecs.ListTaskDefinitionsPaginatorOptions)) EcsListTaskDefinitionsPaginator
	NewListTaskDefinitionFamiliesPaginator(param *ecs.ListTaskDefinitionFamiliesInput, optFns ...func(options *ecs.ListTaskDefinitionFamiliesPaginatorOptions)) EcsListTaskDefinitionFamiliesPaginator
}

type ecsPaginators struct {
	client *ecs.Client
}

func (e *ecsPaginators) NewListClustersPaginator(params *ecs.ListClustersInput, optFns ...func(*ecs.ListClustersPaginatorOptions)) EcsListClustersPaginator {
	return ecs.NewListClustersPaginator(e.client, params, optFns...)
}

type EcsListClustersPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
}

func (e *ecsPaginators) NewListServicesPaginator(params *ecs.ListServicesInput, optFns ...func(*ecs.ListServicesPaginatorOptions)) EcsListServicesPaginator {
	return ecs.NewListServicesPaginator(e.client, params, optFns...)
}

type EcsListServicesPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
}

func (e *ecsPaginators) NewListTaskDefinitionsPaginator(params *ecs.ListTaskDefinitionsInput, optFns ...func(*ecs.ListTaskDefinitionsPaginatorOptions)) EcsListTaskDefinitionsPaginator {
	return ecs.NewListTaskDefinitionsPaginator(e.client, params, optFns...)
}

type EcsListTaskDefinitionsPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error)
}

func (e *ecsPaginators) NewListTaskDefinitionFamiliesPaginator(params *ecs.ListTaskDefinitionFamiliesInput, optFns ...func(*ecs.ListTaskDefinitionFamiliesPaginatorOptions)) EcsListTaskDefinitionFamiliesPaginator {
	return ecs.NewListTaskDefinitionFamiliesPaginator(e.client, params, optFns...)
}

type EcsListTaskDefinitionFamiliesPaginator interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (*ecs.ListTaskDefinitionFamiliesOutput, error)
}
