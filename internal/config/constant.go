package config

import "os"

var (
	SystemName = os.Getenv("USER_SERVICE_SYSTEM_NAME")
	ServiceEnv = os.Getenv("USER_SERVICE_SERVICE_ENV")
	LogLevel   = os.Getenv("USER_SERVICE_LOG_LEVEL")
)

const (
	RoleDeveloper = "developer"
	RoleAdmin     = "admin"

	JWTTokenSalt  = "sad342mslfd23412sdfsdf1234hgf"
	JWTSecret     = "$3cr3t"
	IsoTimeLayout = "2006-01-02T15:04:05Z" // Формат ISO 8601

	ParamID    = "id"
	ParamRoles = "roles"
)
