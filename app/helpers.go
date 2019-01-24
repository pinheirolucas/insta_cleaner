package app

import (
	"context"
	"net/http"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
)

var (
	post = method(http.MethodPost)
	get  = method(http.MethodGet)
)

func method(m string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getSession(ctx context.Context, app *firebase.App, store string) (*storage.ObjectHandle, error) {
	scli, err := app.Storage(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "(*firebase.App).Storage")
	}

	bucket, err := scli.DefaultBucket()
	if err != nil {
		return nil, errors.Wrap(err, "(*storage.Client).DefaultBucket")
	}

	return bucket.Object(store), nil
}
