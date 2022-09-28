# AWS Elastic Container Registry (ECR) cleaner

Removes old, unused images from AWS Elastic Container Registry (ECR).

## Why have we built it?

The number and size of images in your docker registry grow over time and can cost a significant amount of money. It is a
good practice to have a cleaning policy of unused, old images. There are ready-to-use solutions like ECR Lifecycle
policies, but they lack some features - mainly they **don't check if an images is still in use**. We consider using them
**dangerous** - for example if your ECS Service can't find an image when scaling out, it would fail to start your app.
That's why we decided to write a Lambda function that will check periodically for old, unused images.

## How does it work?

The first step is scanning for images that are currently in use and putting the set in memory:

- listing all ECS Services in all ECS Clusters and taking image ids from task definitions,
- listing all Lambda functions with package type "Image" and taking image ids from them,
- listing all App Runner services and taking their image ids.

The second step is iterating over all images in ECR repositories tagged with `BoxCleanerEnabled` set to `true` and for
every image checking if:

- it is older than a threshold (default 30 days),
- it is unused (not present in the set of images that are currently in use).

Any unused, old image is being removed if the `DRY_RUN` environment variable is set to `false`. If you don't set
the `DRY_RUN` environment variable or the value is different than `false`, ECR cleaner will only put line in the logs.

**We strongly advise you to start with `DRY_RUN=true`!**

## Limitations

Only ECS, Lambda, and App Runner are supported. ECR cleaner will not check for any images used by any other service -
for example **EKS is currently not supported**. Also, we **do not check for containers used in different AWS accounts or
different regions**.

There can be some problems with RAM if you use a lot of images as we store them in memory. But we
use `map[string]struct{}` to mitigate the risk.

## Usage

ECR cleaner is designed to be used as an AWS Lambda container image. The Lambda can be triggered by AWS EventBridge
Schedule (we use `cron(0 0 * * ? *)`)

You can use our docker image `ghcr.io/devopsbox-io/aws-ecr-cleaner:v0.1.0` (replace the tag with another version when
appropriate) or build your own image downloading the binary in your Dockerfile (Set the `ECR_CLEANER_SHA256`
and `ECR_CLEANER_VERSION` variables appropriately):

```dockerfile
FROM ubuntu:22.04

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

ENV ECR_CLEANER_SHA256=2c713721af30c4c9380324816bd122f469f0780abc9f86fff62e375d45c61272 \
    ECR_CLEANER_VERSION=0.1.0
RUN curl -L https://github.com/devopsbox-io/aws-ecr-cleaner/releases/download/v${ECR_CLEANER_VERSION}/aws-ecr-cleaner-${ECR_CLEANER_VERSION}-linux-amd64 \
        -o /usr/local/bin/aws-ecr-cleaner && \
    echo "${ECR_CLEANER_SHA256} /usr/local/bin/aws-ecr-cleaner" | sha256sum --check && \
    chmod +x /usr/local/bin/aws-ecr-cleaner

CMD [ "/usr/local/bin/aws-ecr-cleaner" ]
```

Either way, **you have to push ECR cleaner image to an AWS ECR repository in your AWS account before using it in
Lambda**.

The Lambda IAM execution role will need a usual `AWSLambdaBasicExecutionRole` policy and additionally a policy with the
following permissions:

```json
{
  "Statement": [
    {
      "Action": [
        "apprunner:DescribeService",
        "apprunner:ListServices",
        "ecr:BatchDeleteImage",
        "ecr:DescribeImages",
        "ecr:DescribeRepositories",
        "ecr:ListTagsForResource",
        "ecs:DescribeServices",
        "ecs:DescribeTaskDefinition",
        "ecs:ListClusters",
        "ecs:ListServices",
        "lambda:GetFunction",
        "lambda:ListFunctions"
      ],
      "Effect": "Allow",
      "Resource": [
        "*"
      ]
    }
  ],
  "Version": "2012-10-17"
}
```

### Settings

#### Environment variables

- `DEFAULT_KEEP_DAYS` - integer in days, default `30`; ECR cleaner will not remove images younger than value of this
  environment variable
- `DRY_RUN` - boolean, default `true`; if set to `false`, ECR cleaner will start removing images, any other value means
  that ECR cleaner will only put a `Found unused image, should be removed` line to the logs

#### Repository tags

- `BoxCleanerEnabled` - boolean; only repositories with this tag set to `true` will be cleaned
- `BoxCleanerKeepDays` - integer in days; you can override the `DEFAULT_KEEP_DAYS` for each repository using this tag

## Known issues

If you have a lot of old images and there are throttling errors (`error ThrottlingException: Rate exceeded`), just rerun
the process - it is perfectly normal.

## Useful commands related to development

### Generating mocks

```shell
mockgen -source=internal/pkg/aws/apprunner.go -destination=internal/pkg/aws/apprunner_mock.go -package=aws
mockgen -source=internal/pkg/aws/ecr.go -destination=internal/pkg/aws/ecr_mock.go -package=aws
mockgen -source=internal/pkg/aws/ecs.go -destination=internal/pkg/aws/ecs_mock.go -package=aws
mockgen -source=internal/pkg/aws/lambda.go -destination=internal/pkg/aws/lambda_mock.go -package=aws
```

### Running unit tests with coverage

```shell
mkdir -p build/test-results \
  && go test -coverpkg=./... -coverprofile=build/test-results/coverage.out ./... \
  && go tool cover -html=build/test-results/coverage.out -o build/test-results/coverage.html
```
