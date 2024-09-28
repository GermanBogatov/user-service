package config

import "os"

var (
	SystemName = os.Getenv("USER_SERVICE_SYSTEM_NAME")
	ServiceEnv = os.Getenv("USER_SERVICE_SERVICE_ENV")
	LogLevel   = os.Getenv("USER_SERVICE_LOG_LEVEL")
)

const (
	XUserID = "X-USER-ID"
)
