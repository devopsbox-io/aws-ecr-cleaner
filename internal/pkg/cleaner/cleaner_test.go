package cleaner

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/ptr"
	boxaws "github.com/devopsbox-io/aws-ecr-cleaner/internal/pkg/aws"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestCleaner(t *testing.T) {
	t.Parallel()

	type testData struct {
		config         Config
		usedImgs       map[string]struct{}
		existingImages [][]repositoryData
	}

	type deleteImageData struct {
		repositoryName string
		dockerTag      *string
		digest         *string
	}

	tests := map[string]struct {
		input    testData
		expected []deleteImageData
	}{
		"Found old image": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{
				{
					repositoryName: "repo1",
					dockerTag:      ptr.String("v1"),
				},
			},
		},
		"Found young image": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-02T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
		"Found used image": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{
					"repo1uri:v1": {},
				},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
		"Dry run": {
			input: testData{
				config: Config{
					DryRun:          true,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
		"Box cleaner disabled": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "false",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
		"No box cleaner tag": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
		"Non default keep days": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled":  "true",
								"BoxCleanerKeepDays": "29",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-02T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{
				{
					repositoryName: "repo1",
					dockerTag:      ptr.String("v1"),
				},
			},
		},
		"Multiple images": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
						{
							name: "repo2",
							uri:  "repo2uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
					{
						{
							name: "repo3",
							uri:  "repo3uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
						{
							name: "repo4",
							uri:  "repo4uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
								{
									{
										digest: "v2Digest",
										dockerTags: []string{
											"v2",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
									{
										digest: "v3v4Digest",
										dockerTags: []string{
											"v3",
											"v4",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{
				{
					repositoryName: "repo1",
					dockerTag:      ptr.String("v1"),
				},
				{
					repositoryName: "repo2",
					dockerTag:      ptr.String("v1"),
				},
				{
					repositoryName: "repo3",
					dockerTag:      ptr.String("v1"),
				},
				{
					repositoryName: "repo4",
					dockerTag:      ptr.String("v1"),
				},
				{
					repositoryName: "repo4",
					dockerTag:      ptr.String("v2"),
				},
				{
					repositoryName: "repo4",
					dockerTag:      ptr.String("v3"),
				},
				{
					repositoryName: "repo4",
					dockerTag:      ptr.String("v4"),
				},
			},
		},
		"Image without tag": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest:        "imageDigest",
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{
				{
					repositoryName: "repo1",
					digest:         ptr.String("imageDigest"),
				},
			},
		},
		"Found untagged used image by digest": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{
					"repo1uri@imageDigest": {},
				},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest:        "imageDigest",
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
		"Found tagged used image by digest": {
			input: testData{
				config: Config{
					DryRun:          false,
					DefaultKeepDays: 30,
				},
				usedImgs: map[string]struct{}{
					"repo1uri@v1Digest": {},
				},
				existingImages: [][]repositoryData{
					{
						{
							name: "repo1",
							uri:  "repo1uri",
							tags: map[string]string{
								"BoxCleanerEnabled": "true",
							},
							images: [][]imageData{
								{
									{
										digest: "v1Digest",
										dockerTags: []string{
											"v1",
										},
										imagePushedAt: testTimeParse(t, "2022-08-01T00:00:00Z"),
									},
								},
							},
						},
					},
				},
			},
			expected: []deleteImageData{},
		},
	}

	for name, testCase := range tests {
		// capture range variables
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockAwsProvider := boxaws.NewMockProvider(ctrl)
			startTime := testTimeParse(t, "2022-08-31T00:00:01Z")

			mockUsedImages(ctrl, mockAwsProvider, testCase.input.usedImgs)

			mockExistingImages(ctrl, mockAwsProvider, testCase.input.existingImages)

			for _, expectedDelete := range testCase.expected {
				mockAwsProvider.MockEcrClient.EXPECT().BatchDeleteImage(gomock.Any(), &ecr.BatchDeleteImageInput{
					ImageIds: []types.ImageIdentifier{
						{
							ImageTag:    expectedDelete.dockerTag,
							ImageDigest: expectedDelete.digest,
						},
					},
					RepositoryName: aws.String(expectedDelete.repositoryName),
				}).Return(&ecr.BatchDeleteImageOutput{}, nil)
			}

			err := (&Cleaner{
				awsProvider: mockAwsProvider.Provider,
				config:      testCase.input.config,
			}).Clean(startTime)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type imageData struct {
	digest        string
	dockerTags    []string
	imagePushedAt time.Time
}

type repositoryData struct {
	name   string
	uri    string
	tags   map[string]string
	images [][]imageData
}

func mockExistingImages(ctrl *gomock.Controller, mockAwsProvider *boxaws.MockProvider, existingImages [][]repositoryData) {
	mockEcrDescribeResourcesPaginator := boxaws.NewMockEcrDescribeRepositoriesPaginator(ctrl)
	mockAwsProvider.MockEcrPaginators.EXPECT().NewDescribeRepositoriesPaginator(gomock.Any()).Return(mockEcrDescribeResourcesPaginator)
	for _, repositories := range existingImages {
		mockEcrDescribeResourcesPaginator.EXPECT().HasMorePages().Return(true)

		ecrRepositories := make([]types.Repository, len(repositories))
		for i, repoData := range repositories {
			ecrRepositories[i] = types.Repository{
				RepositoryName: aws.String(repoData.name),
				RepositoryArn:  aws.String(fmt.Sprintf("%vArn", repoData.name)),
				RepositoryUri:  aws.String(repoData.uri),
			}
		}

		mockEcrDescribeResourcesPaginator.EXPECT().NextPage(gomock.Any()).Return(&ecr.DescribeRepositoriesOutput{
			Repositories: ecrRepositories,
		}, nil)

		for i, repoData := range repositories {
			tags := make([]types.Tag, 0, len(repoData.tags))
			for key, value := range repoData.tags {
				tags = append(tags, types.Tag{
					Key:   aws.String(key),
					Value: aws.String(value),
				})
			}

			mockAwsProvider.MockEcrClient.EXPECT().ListTagsForResource(gomock.Any(), &ecr.ListTagsForResourceInput{
				ResourceArn: ecrRepositories[i].RepositoryArn,
			}).Return(&ecr.ListTagsForResourceOutput{
				Tags: tags,
			}, nil)

			if tag, ok := repoData.tags["BoxCleanerEnabled"]; ok && tag == "true" {
				mockDescribeImagesPaginator := boxaws.NewMockEcrDescribeImagesPaginator(ctrl)
				mockAwsProvider.MockEcrPaginators.EXPECT().NewDescribeImagesPaginator(&ecr.DescribeImagesInput{
					RepositoryName: aws.String(repoData.name),
				}).Return(mockDescribeImagesPaginator)

				for _, imagesPage := range repoData.images {
					mockDescribeImagesPaginator.EXPECT().HasMorePages().Return(true)

					ecrImages := make([]types.ImageDetail, len(imagesPage))
					for j, image := range imagesPage {
						ecrImages[j] = types.ImageDetail{
							ImagePushedAt:  aws.Time(image.imagePushedAt),
							ImageDigest:    aws.String(image.digest),
							ImageTags:      image.dockerTags,
							RepositoryName: aws.String(repoData.name),
						}
					}

					mockDescribeImagesPaginator.EXPECT().NextPage(gomock.Any()).Return(&ecr.DescribeImagesOutput{
						ImageDetails: ecrImages,
					}, nil)
				}
				mockDescribeImagesPaginator.EXPECT().HasMorePages().Return(false)
			}
		}
	}
	mockEcrDescribeResourcesPaginator.EXPECT().HasMorePages().Return(false)
}

func testTimeParse(t *testing.T, timeStr string) time.Time {
	startTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Fatal(err)
	}
	return startTime
}

func mockUsedImages(ctrl *gomock.Controller, mockAwsProvider *boxaws.MockProvider, used map[string]struct{}) {
	services := make(map[string]string, len(used))

	for usedImage := range used {
		services[fmt.Sprintf("%vService", usedImage)] = usedImage
	}

	mockSsm(ctrl, mockAwsProvider, [][]string{
		{
			"mock-aws-region",
		},
	}, false)

	mockEcs(ctrl, mockAwsProvider, []map[string][]map[string]string{
		{
			"cluster1": {
				services,
			},
		},
	})

	mockLambda(ctrl, mockAwsProvider, []map[string]lambdaMockResult{})

	mockAppRunner(ctrl, mockAwsProvider, []map[string]string{})
}
