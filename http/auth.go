package http

import (
	"fmt"
	"github.com/uxland/ga-go-auth/shared"
	"net/http"
)

func Authenticate(r *http.Request, scheme shared.AuthenticationScheme, context string) {
	r.Header.Set(shared.AuthorizationHeader, fmt.Sprintf("%s %s", scheme, context))
}
