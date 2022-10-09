package cleaner

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	boxaws "github.com/devopsbox-io/aws-ecr-cleaner/internal/pkg/aws"
	gerrors "github.com/pkg/errors"
)

const AppRunnerRegionsSsmParametersPath = "/aws/service/global-infrastructure/services/apprunner/regions"

type usedImages struct {
	awsProvider *boxaws.Provider
}

func (u *usedImages) getImages() (map[string]struct{}, error) {
	imageSet := make(map[string]struct{})

	err := u.getEcsUsedImages(imageSet)
	if err != nil {
		return nil, gerrors.Wrapf(err, "error getting images used by ECS")
	}

	err = u.getLambdaUsedImages(imageSet)
	if err != nil {
		return nil, gerrors.Wrapf(err, "error getting images used by Lambda")
	}

	appRunnerEnabled, err := u.checkAppRunnerEnabledInRegion()
	if err != nil {
		return nil, gerrors.Wrapf(err, "error checking if App Runner is enabled in region")
	}

	if appRunnerEnabled {
		err = u.getAppRunnerUsedImages(imageSet)
		if err != nil {
			return nil, gerrors.Wrapf(err, "error getting images used by App Runner")
		}
	} else {
		logger.Info("App Runner not available in this region", "region", u.awsProvider.Region)
	}

	return imageSet, nil
}

func (u *usedImages) getEcsUsedImages(imageSet map[string]struct{}) error {
	ecsPaginators := u.awsProvider.EcsPaginators
	ecsClient := u.awsProvider.EcsClient

	listClustersPaginator := ecsPaginators.NewListClustersPaginator(&ecs.ListClustersInput{})
	for listClustersPaginator.HasMorePages() {
		listClusterPage, err := listClustersPaginator.NextPage(context.TODO())
		if err != nil {
			return gerrors.Wrapf(err, "cannot get list ECS clusters page")
		}

		for _, clusterArn := range listClusterPage.ClusterArns {

			listServicesPaginator := ecsPaginators.NewListServicesPaginator(&ecs.ListServicesInput{
				Cluster: aws.String(clusterArn),
			})

			for listServicesPaginator.HasMorePages() {
				listServicesPage, err := listServicesPaginator.NextPage(context.TODO())
				if err != nil {
					return gerrors.Wrapf(err, "cannot get list ECS services page")
				}

				describeServicesOutput, err :=
					ecsClient.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{
						Services: listServicesPage.ServiceArns,
						Cluster:  aws.String(clusterArn),
					})
				if err != nil {
					return gerrors.Wrapf(err, "cannot describe ECS services")
				}

				for _, service := range describeServicesOutput.Services {

					describeTaskDefinitionOutput, err :=
						ecsClient.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
							TaskDefinition: service.TaskDefinition,
						})
					if err != nil {
						return gerrors.Wrapf(err, "cannot describe ECS task definitions")
					}

					for _, container := range describeTaskDefinitionOutput.TaskDefinition.ContainerDefinitions {
						image := *container.Image

						logger.Debug("Found image used by ECS service",
							"image", image, "ecsService", *service.ServiceName)

						imageSet[image] = struct{}{}
					}
				}
			}

		}

	}

	return nil
}

func (u *usedImages) getLambdaUsedImages(imageSet map[string]struct{}) error {
	lambdaPaginators := u.awsProvider.LambdaPaginators
	lambdaClient := u.awsProvider.LambdaClient

	listFunctionsPaginator := lambdaPaginators.NewListFunctionsPaginator(&lambda.ListFunctionsInput{})
	for listFunctionsPaginator.HasMorePages() {
		page, err := listFunctionsPaginator.NextPage(context.TODO())
		if err != nil {
			return gerrors.Wrapf(err, "cannot get list Lambda functions page")
		}

		for _, lambdaFunction := range page.Functions {

			if lambdaFunction.PackageType == lambdatypes.PackageTypeImage {

				getFunctionOutput, err := lambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{
					FunctionName: lambdaFunction.FunctionArn,
				})
				if err != nil {
					return gerrors.Wrapf(err, "cannot get Lambda function")
				}

				image := *getFunctionOutput.Code.ImageUri

				logger.Debug("Found image used by Lambda",
					"image", image, "lambda", *lambdaFunction.FunctionName)

				imageSet[image] = struct{}{}
			}
		}
	}

	return nil
}

func (u *usedImages) getAppRunnerUsedImages(imageSet map[string]struct{}) error {
	appRunnerPaginators := u.awsProvider.AppRunnerPaginators
	appRunnerClient := u.awsProvider.AppRunnerClient

	listServicesPaginator := appRunnerPaginators.NewListServicesPaginator(&apprunner.ListServicesInput{})

	for listServicesPaginator.HasMorePages() {
		page, err := listServicesPaginator.NextPage(context.TODO())
		if err != nil {
			return gerrors.Wrapf(err, "cannot get list App Runner services page")
		}

		for _, serviceSummary := range page.ServiceSummaryList {
			describeServiceOutput, err := appRunnerClient.DescribeService(context.TODO(), &apprunner.DescribeServiceInput{
				ServiceArn: serviceSummary.ServiceArn,
			})
			if err != nil {
				return gerrors.Wrapf(err, "cannot describe App Runner service")
			}

			serviceImageRepository := describeServiceOutput.Service.SourceConfiguration.ImageRepository
			if serviceImageRepository != nil {
				image := *serviceImageRepository.ImageIdentifier

				logger.Debug("Found image used by App Runner",
					"image", image, "appRunnerService", *serviceSummary.ServiceName)

				imageSet[image] = struct{}{}
			}
		}

	}

	return nil
}

func (u *usedImages) checkAppRunnerEnabledInRegion() (bool, error) {
	ssmPaginators := u.awsProvider.SsmPaginators
	awsRegion := u.awsProvider.Region

	getParametersByPathPaginator := ssmPaginators.NewGetParametersByPathPaginator(&ssm.GetParametersByPathInput{
		Path: aws.String(AppRunnerRegionsSsmParametersPath),
	})

	for getParametersByPathPaginator.HasMorePages() {
		page, err := getParametersByPathPaginator.NextPage(context.TODO())
		if err != nil {
			return false, gerrors.Wrapf(err, "cannot get get ssm parameters by path page")
		}

		for _, parameter := range page.Parameters {
			supportedRegion := *parameter.Value

			if awsRegion == supportedRegion {
				return true, nil
			}
		}
	}

	return false, nil
}
