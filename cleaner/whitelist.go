package cleaner

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// WhitelistService ...
type WhitelistService interface {
	IsIn(id int64) (bool, error)
}

// FirebaseWhitelistService ...
type FirebaseWhitelistService struct {
	client *firestore.Client
}

// NewFirebaseWhitelistService ...
func NewFirebaseWhitelistService(client *firestore.Client) *FirebaseWhitelistService {
	return &FirebaseWhitelistService{client}
}

// IsIn ...
func (s *FirebaseWhitelistService) IsIn(id int64) (bool, error) {
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

// FakeWhitelistService ...
type FakeWhitelistService func() (bool, error)

// IsIn ...
func (s FakeWhitelistService) IsIn(id int64) (bool, error) {
	return s()
}
