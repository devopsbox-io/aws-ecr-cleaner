package cleaner

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	apprunnertypes "github.com/aws/aws-sdk-go-v2/service/apprunner/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/ptr"
	boxaws "github.com/devopsbox-io/aws-ecr-cleaner/internal/pkg/aws"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetUsedImages(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockAwsProvider := boxaws.NewMockProvider(ctrl)

	mockSsm(ctrl, mockAwsProvider, [][]string{
		{
			"region1",
		},
		{
			"mock-aws-region",
		},
	}, false)

	mockEcs(ctrl, mockAwsProvider, []map[string][]map[string]string{
		{
			"cluster1": {
				{
					"ecsService1": "image1:v1",
				},
				{
					"ecsService2": "image2:v1",
				},
			},
			"cluster2": {
				{
					"ecsService3": "image1:v2",
				},
				{
					"ecsService4": "image2:v2",
				},
			},
		},
		{
			"cluster3": {
				{
					"ecsService5": "duplicatedImage1:v1",
				},
				{
					"ecsService6":  "image3:v1",
					"ecsService12": "duplicatedImage5:v1",
				},
			},
			"cluster4": {
				{
					"ecsService7": "duplicatedImage1:v1",
				},
				{
					"ecsService8": "image3:v2",
				},
				{
					"ecsService9":  "duplicatedImage4:v1",
					"ecsService10": "image6:v1",
					"ecsService11": "duplicatedImage5:v1",
				},
			},
		},
	})

	mockLambda(ctrl, mockAwsProvider, []map[string]lambdaMockResult{
		{
			"lambda1": {
				image:       ptr.String("image4:v1"),
				packageType: lambdatypes.PackageTypeImage,
			},
			"lambda2": {
				image:       ptr.String("duplicatedImage2:v1"),
				packageType: lambdatypes.PackageTypeImage,
			},
			"lambda3": {
				packageType: lambdatypes.PackageTypeZip,
			},
		},
		{
			"lambda4": {
				image:       ptr.String("duplicatedImage2:v1"),
				packageType: lambdatypes.PackageTypeImage,
			},
			"lambda5": {
				image:       ptr.String("image4:v2"),
				packageType: lambdatypes.PackageTypeImage,
			},
			"lambda6": {
				image:       ptr.String("duplicatedImage4:v1"),
				packageType: lambdatypes.PackageTypeImage,
			},
		},
	})

	mockAppRunner(ctrl, mockAwsProvider, []map[string]string{
		{
			"appRunnerService1": "image5:v1",
			"appRunnerService2": "duplicatedImage3:v1",
		},
		{
			"appRunnerService3": "image5:v2",
			"appRunnerService4": "duplicatedImage3:v1",
			"appRunnerService5": "duplicatedImage4:v1",
		},
	})

	expectedImages := map[string]struct{}{
		"image1:v1":           {},
		"image2:v1":           {},
		"image1:v2":           {},
		"image2:v2":           {},
		"duplicatedImage1:v1": {},
		"image3:v1":           {},
		"image3:v2":           {},
		"duplicatedImage4:v1": {},
		"image4:v1":           {},
		"duplicatedImage2:v1": {},
		"image4:v2":           {},
		"image5:v1":           {},
		"duplicatedImage3:v1": {},
		"image5:v2":           {},
		"duplicatedImage5:v1": {},
		"image6:v1":           {},
	}

	images, err := (&usedImages{
		awsProvider: mockAwsProvider.Provider,
	}).getImages()
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(
		expectedImages,
		images,
	)
	if diff != "" {
		t.Error(diff)
	}
}

func TestAppRunnerNotAvailable(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockAwsProvider := boxaws.NewMockProvider(ctrl)

	mockSsm(ctrl, mockAwsProvider, [][]string{
		{
			"region2",
		},
	}, true)

	mockEcs(ctrl, mockAwsProvider, []map[string][]map[string]string{
		{
			"cluster1": {
				{
					"ecsService1": "image1:v1",
				},
			},
		},
	})

	mockLambda(ctrl, mockAwsProvider, []map[string]lambdaMockResult{
		{
			"lambda1": {
				image:       ptr.String("image2:v1"),
				packageType: lambdatypes.PackageTypeImage,
			},
		},
	})

	expectedImages := map[string]struct{}{
		"image1:v1": {},
		"image2:v1": {},
	}

	images, err := (&usedImages{
		awsProvider: mockAwsProvider.Provider,
	}).getImages()
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(
		expectedImages,
		images,
	)
	if diff != "" {
		t.Error(diff)
	}
}

func mockSsm(
	ctrl *gomock.Controller,
	mockAwsProvider *boxaws.MockProvider,
	mockResult [][]string,
	expectSsmGetParametersByPathPaginatorAllPages bool,
) {

	mockSsmGetParametersByPathPaginator := boxaws.NewMockSsmGetParametersByPathPaginator(ctrl)
	mockAwsProvider.MockSsmPaginators.EXPECT().NewGetParametersByPathPaginator(gomock.Any()).Return(mockSsmGetParametersByPathPaginator)

	for _, getParametersByPathPage := range mockResult {
		mockSsmGetParametersByPathPaginator.EXPECT().HasMorePages().Return(true)

		parameters := make([]ssmtypes.Parameter, len(getParametersByPathPage))

		for i, parameterValue := range getParametersByPathPage {
			parameters[i] = ssmtypes.Parameter{
				Value: aws.String(parameterValue),
			}
		}

		mockSsmGetParametersByPathPaginator.EXPECT().NextPage(gomock.Any()).Return(&ssm.GetParametersByPathOutput{
			Parameters: parameters,
		}, nil)
	}
	if expectSsmGetParametersByPathPaginatorAllPages {
		mockSsmGetParametersByPathPaginator.EXPECT().HasMorePages().Return(false)
	}
}

func mockEcs(ctrl *gomock.Controller, mockAwsProvider *boxaws.MockProvider, mockResult []map[string][]map[string]string) {
	mockEcsListClustersPaginator := boxaws.NewMockEcsListClustersPaginator(ctrl)
	mockAwsProvider.MockEcsPaginators.EXPECT().NewListClustersPaginator(gomock.Any()).Return(mockEcsListClustersPaginator)

	for _, listClustersPage := range mockResult {
		mockEcsListClustersPaginator.EXPECT().HasMorePages().Return(true)

		clusterArns := make([]string, 0, len(listClustersPage))
		clusterNames := make([]string, 0, len(listClustersPage))
		for cluster := range listClustersPage {
			clusterArns = append(clusterArns, fmt.Sprintf("%vArn", cluster))
			clusterNames = append(clusterNames, cluster)
		}

		mockEcsListClustersPaginator.EXPECT().NextPage(gomock.Any()).Return(&ecs.ListClustersOutput{
			ClusterArns: clusterArns,
		}, nil)

		for i, clusterName := range clusterNames {
			clusterArn := clusterArns[i]
			listServicesPages := listClustersPage[clusterName]

			mockEcsListServicesPaginator := boxaws.NewMockEcsListServicesPaginator(ctrl)
			mockAwsProvider.MockEcsPaginators.EXPECT().NewListServicesPaginator(&ecs.ListServicesInput{
				Cluster: aws.String(clusterArn),
			}).Return(mockEcsListServicesPaginator)

			for _, listServicesPage := range listServicesPages {
				mockEcsListServicesPaginator.EXPECT().HasMorePages().Return(true)

				serviceArns := make([]string, 0, len(listServicesPage))
				ecsServices := make([]ecstypes.Service, 0, len(listServicesPage))
				for service := range listServicesPage {
					serviceArns = append(serviceArns, fmt.Sprintf("%vArn", service))
					ecsServices = append(ecsServices, ecstypes.Service{
						ServiceName:    aws.String(service),
						TaskDefinition: aws.String(fmt.Sprintf("%vTaskDefinition", service)),
					})
				}

				mockEcsListServicesPaginator.EXPECT().NextPage(gomock.Any()).Return(&ecs.ListServicesOutput{
					ServiceArns: serviceArns,
				}, nil)

				mockAwsProvider.MockEcsClient.EXPECT().DescribeServices(gomock.Any(), &ecs.DescribeServicesInput{
					Services: serviceArns,
					Cluster:  aws.String(clusterArn),
				}).Return(&ecs.DescribeServicesOutput{
					Services: ecsServices,
				}, nil)

				for _, ecsService := range ecsServices {
					image := listServicesPage[*ecsService.ServiceName]

					mockAwsProvider.MockEcsClient.EXPECT().DescribeTaskDefinition(gomock.Any(), &ecs.DescribeTaskDefinitionInput{
						TaskDefinition: ecsService.TaskDefinition,
					}).Return(&ecs.DescribeTaskDefinitionOutput{
						TaskDefinition: &ecstypes.TaskDefinition{
							ContainerDefinitions: []ecstypes.ContainerDefinition{
								{
									Image: aws.String(image),
								},
							},
						},
					}, nil)
				}
			}
			mockEcsListServicesPaginator.EXPECT().HasMorePages().Return(false)
		}
	}
	mockEcsListClustersPaginator.EXPECT().HasMorePages().Return(false)
}

type lambdaMockResult struct {
	image       *string
	packageType lambdatypes.PackageType
}

func mockLambda(ctrl *gomock.Controller, mockAwsProvider *boxaws.MockProvider, mockResult []map[string]lambdaMockResult) {
	mockLambdaListFunctionsPaginator := boxaws.NewMockLambdaListFunctionsPaginator(ctrl)
	mockAwsProvider.MockLambdaPaginators.EXPECT().NewListFunctionsPaginator(gomock.Any()).Return(mockLambdaListFunctionsPaginator)

	for _, listFunctionsPage := range mockResult {
		mockLambdaListFunctionsPaginator.EXPECT().HasMorePages().Return(true)

		lambdaFunctions := make([]lambdatypes.FunctionConfiguration, 0, len(listFunctionsPage))
		for functionName, functionParameters := range listFunctionsPage {
			lambdaFunctions = append(lambdaFunctions, lambdatypes.FunctionConfiguration{
				FunctionName: aws.String(functionName),
				FunctionArn:  aws.String(fmt.Sprintf("%vArn", functionName)),
				PackageType:  functionParameters.packageType,
			})
		}

		mockLambdaListFunctionsPaginator.EXPECT().NextPage(gomock.Any()).Return(&lambda.ListFunctionsOutput{
			Functions: lambdaFunctions,
		}, nil)

		for _, function := range lambdaFunctions {
			functionParameters := listFunctionsPage[*function.FunctionName]

			if functionParameters.packageType == lambdatypes.PackageTypeImage {
				mockAwsProvider.MockLambdaClient.EXPECT().GetFunction(gomock.Any(), &lambda.GetFunctionInput{
					FunctionName: function.FunctionArn,
				}).Return(&lambda.GetFunctionOutput{
					Code: &lambdatypes.FunctionCodeLocation{
						ImageUri: functionParameters.image,
					},
				}, nil)
			}
		}
	}
	mockLambdaListFunctionsPaginator.EXPECT().HasMorePages().Return(false)
}

func mockAppRunner(ctrl *gomock.Controller, mockAwsProvider *boxaws.MockProvider, mockResult []map[string]string) {
	mockApprunnerListServicesPaginator := boxaws.NewMockAppRunnerListServicesPaginator(ctrl)
	mockAwsProvider.MockAppRunnerPaginators.EXPECT().NewListServicesPaginator(gomock.Any()).Return(mockApprunnerListServicesPaginator)

	for _, listServicesPage := range mockResult {
		mockApprunnerListServicesPaginator.EXPECT().HasMorePages().Return(true)

		appRunnerServices := make([]apprunnertypes.ServiceSummary, 0, len(listServicesPage))
		for service := range listServicesPage {
			appRunnerServices = append(appRunnerServices, apprunnertypes.ServiceSummary{
				ServiceArn:  aws.String(fmt.Sprintf("%vArn", service)),
				ServiceName: aws.String(service),
			})
		}

		mockApprunnerListServicesPaginator.EXPECT().NextPage(gomock.Any()).Return(&apprunner.ListServicesOutput{
			ServiceSummaryList: appRunnerServices,
		}, nil)

		for _, service := range appRunnerServices {
			image := listServicesPage[*service.ServiceName]

			mockAwsProvider.MockAppRunnerClient.EXPECT().DescribeService(gomock.Any(), &apprunner.DescribeServiceInput{
				ServiceArn: service.ServiceArn,
			}).Return(&apprunner.DescribeServiceOutput{
				Service: &apprunnertypes.Service{
					SourceConfiguration: &apprunnertypes.SourceConfiguration{
						ImageRepository: &apprunnertypes.ImageRepository{
							ImageIdentifier: aws.String(image),
						},
					},
				},
			}, nil)
		}
	}
	mockApprunnerListServicesPaginator.EXPECT().HasMorePages().Return(false)
}
