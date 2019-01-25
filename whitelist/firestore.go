package whitelist

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// FirestoreService ...
type FirestoreService struct {
	client *firestore.Client
}

// NewFirestoreService ...
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{client}
}

// IsIn ...
func (s *FirestoreService) IsIn(id int64) (bool, error) {
	_, err := s.client.Collection("whitelist").Doc(strconv.Itoa(int(id))).Get(context.Background())
	switch grpc.Code(err) {
	case codes.OK:
		// do nothing
	case codes.NotFound:
		return false, nil
	default:
		return false, errors.Wrapf(err, "not able to Get doc ID: %d", id)
	}

	return true, nil
}

// CreateIfNotExists ...
func (s *FirestoreService) CreateIfNotExists(id int64, username string) error {
	snap, err := s.client.Collection("whitelist").Doc(strconv.Itoa(int(id))).Get(context.Background())
	switch grpc.Code(err) {
	case codes.OK:
		if snap.Exists() {
			return nil
		}
	case codes.NotFound:
		// proceed do create
	default:
		return errors.Wrapf(err, "not able to get %s: %d", username, id)
	}

	_, err = snap.Ref.Create(context.Background(), map[string]interface{}{
		"id":       id,
		"username": username,
	})
	if err != nil {
		return errors.Wrap(err, "snap.Ref.Create")
	}

	return nil
}

// Delete ...
func (s *FirestoreService) Delete(id int64) error {
	_, err := s.client.Collection("whitelist").Doc(strconv.Itoa(int(id))).Delete(context.Background())
	if err != nil {
		return errors.Wrapf(err, "not able to delete: %d", id)
	}

	return nil
}
