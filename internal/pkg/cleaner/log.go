package cleaner

import (
	"github.com/hashicorp/go-hclog"
	"os"
)

const logLevelEnvVar = "BOX_LOG"
const defaultLogLevel = "INFO"

var logger = hclog.New(&hclog.LoggerOptions{
	Name:  "aws-ecr-cleaner",
	Level: hclog.LevelFromString(getLogLevel()),
})

func getLogLevel() string {
	value := os.Getenv(logLevelEnvVar)
	if len(value) == 0 {
		return defaultLogLevel
	}
	return value
}
