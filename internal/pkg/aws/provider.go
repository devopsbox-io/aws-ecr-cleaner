package aws

import gerrors "github.com/pkg/errors"

func NewProvider() (*Provider, error) {
	ecsClient, err := newEcsClient()
	if err != nil {
		return nil, gerrors.Wrapf(err, "error creating AWS ECS client")
	}
	lambdaClient, err := newLambdaClient()
	if err != nil {
		return nil, gerrors.Wrapf(err, "error creating AWS Lambda client")
	}
	appRunnerClient, err := newAppRunnerClient()
	if err != nil {
		return nil, gerrors.Wrapf(err, "error creating AWS App Runner client")
	}
	ecrClient, err := newEcrClient()
	if err != nil {
		return nil, gerrors.Wrapf(err, "error creating AWS ECR client")
	}

	return &Provider{
		EcsClient:           ecsClient,
		EcsPaginators:       &ecsPaginators{client: ecsClient},
		LambdaClient:        lambdaClient,
		LambdaPaginators:    &lambdaPaginators{client: lambdaClient},
		AppRunnerClient:     appRunnerClient,
		AppRunnerPaginators: &appRunnerPaginators{client: appRunnerClient},
		EcrClient:           ecrClient,
		EcrPaginators:       &ecrPaginators{client: ecrClient},
	}, nil
}

type Provider struct {
	EcsClient     EcsClient
	EcsPaginators EcsPaginators

	LambdaClient     LambdaClient
	LambdaPaginators LambdaPaginators

	AppRunnerClient     AppRunnerClient
	AppRunnerPaginators AppRunnerPaginators

	EcrClient     EcrClient
	EcrPaginators EcrPaginators
}
