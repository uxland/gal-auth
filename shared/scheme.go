package shared

type AuthenticationScheme = string

const (
	BearerSchema AuthenticationScheme = "Bearer"
)

const AuthorizationHeader = "authorization"
