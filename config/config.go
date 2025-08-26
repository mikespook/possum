package config

import (
	"os"
)

const (
	Development = "development"
	Production  = "production"
	Test        = "test"
)

var (
	defaultEnv = Development
)

func init() {
	defaultEnv = os.Getenv("POSSUM_ENV")
}

// IsDebug returns true if the application is not running in production mode.
func IsDebug() bool {
	return defaultEnv != Production
}

// IsDev returns true if the application is running in development mode.
func IsDev() bool {
	return defaultEnv == Development
}
