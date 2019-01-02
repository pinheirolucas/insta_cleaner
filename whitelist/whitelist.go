package whitelist

import (
	"context"

	"firebase.google.com/go"
	"github.com/pkg/errors"
)

const (
	serviceTypeFirestore        = "firestore"
	serviceTypeRealtimeDatabase = "realtime"
)

// Service ...
type Service interface {
	IsIn(id int64) (bool, error)
	CreateIfNotExists(id int64, username string) error
}

// NewService ...
func NewService(ctx context.Context, serviceType string, app *firebase.App) (Service, error) {
	switch serviceType {
	case serviceTypeFirestore:
		firestore, err := app.Firestore(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "(*firebase.App).Firestore")
		}

		return NewFirestoreService(firestore), nil
	case serviceTypeRealtimeDatabase:
		database, err := app.Database(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "(*firebase.App).Database")
		}

		return NewRealtimeDatabaseService(database), nil
	}

	return nil, errors.Errorf(`invalid service type "%s"`, serviceType)
}
