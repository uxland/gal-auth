package shared

type AuthenticationScheme = string

const (
	BearerSchema AuthenticationScheme = "Bearer"
	BasicSchema  AuthenticationScheme = "Basic"
)

const AuthorizationHeader = "authorization"
