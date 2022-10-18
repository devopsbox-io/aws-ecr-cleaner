package cleaner

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/ptr"
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

func (c *Cleaner) processSingleImage(
	repository types.Repository,
	image types.ImageDetail,
	usedImagesSet map[string]struct{},
	keepDays int,
	startTime time.Time,
) error {

	imageAgeDays := startTime.Sub(*image.ImagePushedAt).Hours() / 24

	if imageAgeDays > float64(keepDays) {

		imageDigest := *image.ImageDigest

		logger.Debug("Found old image", "repository", *repository.RepositoryUri, "imageAgeDays", imageAgeDays)

		for _, imageTag := range image.ImageTags {
			// capture range variables
			imageTag := imageTag

			err := c.processSingleImageReference(imageReference{
				repositoryUri:  *repository.RepositoryUri,
				repositoryName: *repository.RepositoryName,
				digest:         imageDigest,
				tag:            &imageTag,
			}, usedImagesSet)
			if err != nil {
				return gerrors.Wrapf(err, "error processig %v image tag in repository %v", imageTag, *repository.RepositoryName)
			}
		}

		if len(image.ImageTags) == 0 {
			logger.Debug("Found untagged image", "imageDigest", imageDigest)

			err := c.processSingleImageReference(imageReference{
				repositoryUri:  *repository.RepositoryUri,
				repositoryName: *repository.RepositoryName,
				digest:         imageDigest,
			}, usedImagesSet)
			if err != nil {
				return gerrors.Wrapf(err, "error processig %v image tag in repository %v", imageDigest, *repository.RepositoryName)
			}
		}
	}
	return nil
}

type imageReference struct {
	repositoryUri  string
	repositoryName string
	digest         string
	tag            *string
}

func (i imageReference) String() string {
	if i.tag != nil {
		return *i.tagId()
	} else {
		return i.digestId()
	}
}

func (i imageReference) tagId() *string {
	if i.tag != nil {
		return ptr.String(fmt.Sprintf("%v:%v", i.repositoryUri, *i.tag))
	} else {
		return nil
	}
}

func (i imageReference) digestId() string {
	return fmt.Sprintf("%v@%v", i.repositoryUri, i.digest)
}

func (c *Cleaner) processSingleImageReference(reference imageReference, usedImagesSet map[string]struct{}) error {

	if !isImageInUse(reference, usedImagesSet) {
		if c.config.DryRun {
			logger.Info("Found unused image, should be removed",
				"imageReference", reference)
		} else {
			logger.Info("Found unused image, removing",
				"imageReference", reference)

			err := c.deleteImage(reference)
			if err != nil {
				return gerrors.Wrapf(err, "error deleting image %v", reference)
			}
		}
	}
	return nil
}

func isImageInUse(reference imageReference, usedImagesSet map[string]struct{}) bool {
	imageTagId := reference.tagId()
	imageDigestId := reference.digestId()

	if imageTagId != nil {
		if _, ok := usedImagesSet[*imageTagId]; ok {
			logger.Info("Found old image in use", "imageId", *imageTagId)
			return true
		}
	}

	if _, ok := usedImagesSet[imageDigestId]; ok {
		logger.Info("Found old image in use", "imageId", imageDigestId)
		return true
	}

	return false
}

func (c *Cleaner) deleteImage(reference imageReference) error {
	ecrClient := c.awsProvider.EcrClient

	var imageIdentifier types.ImageIdentifier
	if reference.tag != nil {
		imageIdentifier = types.ImageIdentifier{
			ImageTag: reference.tag,
		}
	} else {
		imageIdentifier = types.ImageIdentifier{
			ImageDigest: aws.String(reference.digest),
		}
	}

	_, err := ecrClient.BatchDeleteImage(context.TODO(), &ecr.BatchDeleteImageInput{
		ImageIds: []types.ImageIdentifier{
			imageIdentifier,
		},
		RepositoryName: aws.String(reference.repositoryName),
	})
	if err != nil {
		return gerrors.Wrapf(err, "cannot remove image %v from repository", reference)
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
