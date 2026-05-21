package config

const (
	ENV_FILE          = ".env"
	ENV_FILE_DATABASE = ".env.postgres"

	ENV_APPLICATION_PORT = "APPLICATION_PORT"
	ENV_TOKEN_KEY        = "TOKEN_KEY"
)

const (
	DATABASE_NAME = "nyasama"
)

const (
	BCRYPT_COST       = 10
	SALT_LENGTH       = 16
	PASSWORD_MAX_INDEX = 55
	AUTHENTICATION_TYPE = "Bearer "
)

const (
	CONTEXT_KEY_CLAIMS = "claims"
)

const (
	HTTP_HEADER_AUTHORIZATION = "Authorization"
)

type PageType string

const (
	COMIC   PageType = "comic"
	ARTICLE PageType = "article"
)

const (
	API_REGISTER      = "/register"
	API_LOGIN         = "/login"
	API_GET_USER_INFO = "/user/info"
)
