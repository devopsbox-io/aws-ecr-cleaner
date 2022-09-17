package cleaner

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	boxaws "github.com/devopsbox-io/aws-ecr-cleaner/internal/pkg/aws"
	gerrors "github.com/pkg/errors"
	"strconv"
	"time"
)

type Config struct {
	DryRun          bool
	DefaultKeepDays int
}

func New(awsProvider *boxaws.Provider, config Config) *Cleaner {
	return &Cleaner{
		awsProvider: awsProvider,
		config:      config,
	}
}

type Cleaner struct {
	awsProvider *boxaws.Provider
	config      Config
}

const (
	BoxCleanerEnabledTag  = "BoxCleanerEnabled"
	BoxCleanerKeepDaysTag = "BoxCleanerKeepDays"
)

func (c *Cleaner) Clean(startTime time.Time) error {
	usedImagesSet, err := (&usedImages{awsProvider: c.awsProvider}).getImages()
	if err != nil {
		return gerrors.Wrapf(err, "error getting used images")
	}

	logger.Info("Found used images", "len(usedImagesSet)", len(usedImagesSet))

	ecrPaginators := c.awsProvider.EcrPaginators

	describeRepositoriesPaginator := ecrPaginators.NewDescribeRepositoriesPaginator(&ecr.DescribeRepositoriesInput{})
	for describeRepositoriesPaginator.HasMorePages() {
		describeRepositoriesPage, err := describeRepositoriesPaginator.NextPage(context.TODO())
		if err != nil {
			return gerrors.Wrapf(err, "cannot get describe repositories page")
		}

		for _, repository := range describeRepositoriesPage.Repositories {

			err := c.processSingleRepository(repository, usedImagesSet, startTime)
			if err != nil {
				return gerrors.Wrapf(err, "error processig %v repository", *repository.RepositoryName)
			}
		}
	}

	return nil
}

func (c *Cleaner) processSingleRepository(repository types.Repository, usedImagesSet map[string]struct{}, startTime time.Time) error {
	ecrClient := c.awsProvider.EcrClient

	listTagsForResourceOutput, err := ecrClient.ListTagsForResource(context.TODO(), &ecr.ListTagsForResourceInput{
		ResourceArn: repository.RepositoryArn,
	})
	if err != nil {
		return gerrors.Wrapf(err, "cannot list tags for repository %v", *repository.RepositoryArn)
	}
	repositoryTagsMap := convertTagsToMap(listTagsForResourceOutput.Tags)

	logger.Debug("Found repository tags", "repository", *repository.RepositoryArn, "repositoryTagsMap", repositoryTagsMap)

	if boxCleanerEnabledTagValue, ok := repositoryTagsMap[BoxCleanerEnabledTag]; ok && boxCleanerEnabledTagValue == "true" {

		keepDays := c.countKeepDays(repositoryTagsMap)

		err := c.cleanSingleRepository(repository, usedImagesSet, keepDays, startTime)
		if err != nil {
			return gerrors.Wrapf(err, "error cleaning %v repository", *repository.RepositoryName)
		}
	}
	return nil
}

func (c *Cleaner) cleanSingleRepository(repository types.Repository, usedImagesSet map[string]struct{}, keepDays int, startTime time.Time) error {
	ecrPaginators := c.awsProvider.EcrPaginators

	describeImagesPaginator := ecrPaginators.NewDescribeImagesPaginator(&ecr.DescribeImagesInput{
		RepositoryName: repository.RepositoryName,
	})

	for describeImagesPaginator.HasMorePages() {
		describeImagesPage, err := describeImagesPaginator.NextPage(context.TODO())
		if err != nil {
			return gerrors.Wrapf(err, "cannot get describe images page")
		}
		for _, image := range describeImagesPage.ImageDetails {
			err := c.processSingleImage(repository, image, usedImagesSet, keepDays, startTime)
			if err != nil {
				return gerrors.Wrapf(err, "error processig %v image in repository %v", *image.ImageDigest, *repository.RepositoryName)
			}
		}
	}
	return nil
}

func (c *Cleaner) processSingleImage(repository types.Repository, image types.ImageDetail, usedImagesSet map[string]struct{}, keepDays int, startTime time.Time) error {
	imageAgeDays := startTime.Sub(*image.ImagePushedAt).Hours() / 24

	if imageAgeDays > float64(keepDays) {

		for _, imageTag := range image.ImageTags {
			err := c.processSingleImageTag(repository, imageTag, usedImagesSet, imageAgeDays)
			if err != nil {
				return gerrors.Wrapf(err, "error processig %v image tag in repository %v", imageTag, *repository.RepositoryName)
			}
		}
	}
	return nil
}

func (c *Cleaner) processSingleImageTag(repository types.Repository, imageTag string, usedImagesSet map[string]struct{}, imageAgeDays float64) error {

	imageId := fmt.Sprintf("%v:%v", *repository.RepositoryUri, imageTag)

	if _, ok := usedImagesSet[imageId]; ok {
		logger.Info("Found old used image", "imageId", imageId, "imageAgeDays", imageAgeDays)
	} else {
		if c.config.DryRun {
			logger.Info("Found unused image, should be removed",
				"imageId", imageId, "imageAgeDays", imageAgeDays)
		} else {
			logger.Info("Found unused image, removing",
				"imageId", imageId, "imageAgeDays", imageAgeDays)

			err := c.deleteImage(repository, imageTag)
			if err != nil {
				return gerrors.Wrapf(err, "error deleting image %v", imageId)
			}
		}
	}
	return nil
}

func (c *Cleaner) deleteImage(repository types.Repository, imageTag string) error {
	ecrClient := c.awsProvider.EcrClient

	_, err := ecrClient.BatchDeleteImage(context.TODO(), &ecr.BatchDeleteImageInput{
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: aws.String(imageTag),
			},
		},
		RepositoryName: repository.RepositoryName,
	})
	if err != nil {
		return gerrors.Wrapf(err, "cannot remove image %v from repository %v", imageTag, *repository.RepositoryName)
	}
	return nil
}

func (c *Cleaner) countKeepDays(repositoryTagsMap map[string]string) int {
	keepDays := c.config.DefaultKeepDays
	if boxCleanerKeepDaysTagValueStr, ok := repositoryTagsMap[BoxCleanerKeepDaysTag]; ok {
		boxCleanerKeepDaysTagValue, err := strconv.Atoi(boxCleanerKeepDaysTagValueStr)
		if err == nil {
			keepDays = boxCleanerKeepDaysTagValue
		}
	}
	return keepDays
}

func convertTagsToMap(tags []types.Tag) map[string]string {
	tagsMap := make(map[string]string, len(tags))

	for _, tag := range tags {
		tagsMap[*tag.Key] = *tag.Value
	}

	return tagsMap
}
