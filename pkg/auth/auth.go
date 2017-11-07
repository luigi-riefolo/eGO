package auth

import (
	"context"

	"github.com/luigi-riefolo/alfa/pkg/http"
)

// Token gets the auth token from the context.
func Token(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(http.ContextKey("auth-token")).(string)
	return val, ok
}
