package shared

type UnauthorizedError struct {
}

func (u UnauthorizedError) Error() string {
	return "Unauthorized"
}
