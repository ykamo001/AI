package request

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

var headerKeys = []RequestKey{Origin, ContentType, Cookie, Authorization, XSignature}

func WithRequestHeaders(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for _, key := range headerKeys {
			val := r.Header.Get(string(key))
			if val != "" {
				ctx = context.WithValue(ctx, key, val)
			}
		}

		id, err := uuid.NewRandom()
		if err == nil {
			ctx = context.WithValue(ctx, "id", id.ID())
		}
		
		r = r.WithContext(ctx)
		base.ServeHTTP(w, r)
	})
}
