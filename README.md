

Does not support:
- EKS
- Mutli-account
- Multi-region

## Generating mocks

```shell
mockgen -source=internal/pkg/aws/apprunner.go -destination=internal/pkg/aws/apprunner_mock.go -package=aws
mockgen -source=internal/pkg/aws/ecr.go -destination=internal/pkg/aws/ecr_mock.go -package=aws
mockgen -source=internal/pkg/aws/ecs.go -destination=internal/pkg/aws/ecs_mock.go -package=aws
mockgen -source=internal/pkg/aws/lambda.go -destination=internal/pkg/aws/lambda_mock.go -package=aws
```

## Running unit tests with coverage

```shell
mkdir -p build/test-results \
  && go test -coverpkg=./... -coverprofile=build/test-results/coverage.out ./... \
  && go tool cover -html=build/test-results/coverage.out -o build/test-results/coverage.html
```