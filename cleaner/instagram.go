package cleaner

import (
	"github.com/ahmdrz/goinsta"
)

// InstagramService ...
type InstagramService interface {
	GetFollowingUsers() *goinsta.Users
}

// GoinstaInstagramService ...
type GoinstaInstagramService struct {
	insta *goinsta.Instagram
}

// NewGoinstaInstagramService ...
func NewGoinstaInstagramService(insta *goinsta.Instagram) *GoinstaInstagramService {
	return &GoinstaInstagramService{insta}
}

// GetFollowingUsers ...
func (i *GoinstaInstagramService) GetFollowingUsers() *goinsta.Users {
	return i.insta.Account.Following()
}

// FakeInstagramService ...
type FakeInstagramService func() *goinsta.Users

// GetFollowingUsers ...
func (i FakeInstagramService) GetFollowingUsers() *goinsta.Users {
	return i()
}
