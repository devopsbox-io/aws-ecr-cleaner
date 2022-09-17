package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devopsbox-io/aws-ecr-cleaner/internal/pkg/aws"
	"github.com/devopsbox-io/aws-ecr-cleaner/internal/pkg/cleaner"
	"os"
	"strconv"
	"time"
)

const DefaultKeepDays = 30

func main() {
	awsProvider, err := aws.NewProvider()
	if err != nil {
		panic(err)
	}

	cleanerObj := cleaner.New(awsProvider, cleaner.Config{
		DryRun:          getDryRun(os.LookupEnv),
		DefaultKeepDays: getDefaultKeepDays(os.LookupEnv),
	})

	if isLambda(os.LookupEnv) {
		lambda.Start(func() error {
			return cleanerObj.Clean(time.Now())
		})
	} else {
		err := cleanerObj.Clean(time.Now())
		if err != nil {
			panic(err)
		}
	}
}

func isLambda(lookupEnv func(key string) (string, bool)) bool {
	_, result := lookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	return result
}

func getDefaultKeepDays(lookupEnv func(key string) (string, bool)) int {
	defaultKeepDays := DefaultKeepDays
	defaultKeepDaysStr, isDefaultKeepDaysSet := lookupEnv("DEFAULT_KEEP_DAYS")
	if isDefaultKeepDaysSet {
		parsedDefaultKeepDays, err := strconv.Atoi(defaultKeepDaysStr)
		if err == nil {
			defaultKeepDays = parsedDefaultKeepDays
		}
	}
	return defaultKeepDays
}

func getDryRun(lookupEnv func(key string) (string, bool)) bool {
	dryRun := true
	dryRunStr, isDryRunSet := lookupEnv("DRY_RUN")
	if isDryRunSet && dryRunStr == "false" {
		dryRun = false
	}
	return dryRun
}
