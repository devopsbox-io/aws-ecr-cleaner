package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	gerrors "github.com/pkg/errors"
)

func NewProvider() (*Provider, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, gerrors.Wrapf(err, "cannot load aws config")
	}

	ecsClient := newEcsClient(cfg)
	lambdaClient := newLambdaClient(cfg)
	appRunnerClient := newAppRunnerClient(cfg)
	ecrClient := newEcrClient(cfg)
	ssmClient := newSsmClient(cfg)

	return &Provider{
		Region: cfg.Region,

		EcsClient:           ecsClient,
		EcsPaginators:       &ecsPaginators{client: ecsClient},
		LambdaClient:        lambdaClient,
		LambdaPaginators:    &lambdaPaginators{client: lambdaClient},
		AppRunnerClient:     appRunnerClient,
		AppRunnerPaginators: &appRunnerPaginators{client: appRunnerClient},
		EcrClient:           ecrClient,
		EcrPaginators:       &ecrPaginators{client: ecrClient},
		SsmPaginators:       &ssmPaginators{client: ssmClient},
	}, nil
}

type Provider struct {
	Region string

	EcsClient     EcsClient
	EcsPaginators EcsPaginators

	LambdaClient     LambdaClient
	LambdaPaginators LambdaPaginators

	AppRunnerClient     AppRunnerClient
	AppRunnerPaginators AppRunnerPaginators

	EcrClient     EcrClient
	EcrPaginators EcrPaginators

	SsmPaginators SsmPaginators
}
