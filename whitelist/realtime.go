package whitelist

import (
	"context"
	"fmt"

	"firebase.google.com/go/db"
	"github.com/pkg/errors"
)

// RealtimeDatabaseService ...
type RealtimeDatabaseService struct {
	client *db.Client
}

// NewRealtimeDatabaseService ...
func NewRealtimeDatabaseService(client *db.Client) *RealtimeDatabaseService {
	return &RealtimeDatabaseService{client}
}

// IsIn ...
func (s *RealtimeDatabaseService) IsIn(id int64) (bool, error) {
	path := fmt.Sprintf("whitelist/%d", id)

	result := new(whitelistMember)
	if err := s.client.NewRef(path).Get(context.Background(), result); err != nil {
		return false, errors.Wrapf(err, "not able to get path %s", path)
	}

	return result.ID != 0, nil
}

// CreateIfNotExists ...
func (s *RealtimeDatabaseService) CreateIfNotExists(id int64, username string) error {
	path := fmt.Sprintf("whitelist/%d", id)

	err := s.client.NewRef(path).Set(context.Background(), &whitelistMember{
		ID:       id,
		Username: username,
	})
	if err != nil {
		return errors.Wrapf(err, "not able to set path %s with username: %s and id %d", path, username, id)
	}

	return nil
}

// Delete ...
func (s *RealtimeDatabaseService) Delete(id int64) error {
	path := fmt.Sprintf("whitelist/%d", id)

	if err := s.client.NewRef(path).Delete(context.Background()); err != nil {
		return errors.Wrapf(err, "not able to delete path %s with id %d", path, id)
	}

	return nil
}

type whitelistMember struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}
