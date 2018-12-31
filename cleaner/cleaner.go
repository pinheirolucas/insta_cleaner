package cleaner

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pinheirolucas/insta_cleaner/logger"
)

// Service ...
type Service struct {
	l logger.Logger

	unfollowCooldown time.Duration
	maxUnfollows     uint32

	instagramService InstagramService
	whitelistService WhitelistService
}

// NewService ...
func NewService(instagramService InstagramService, whitelistService WhitelistService, l logger.Logger, options ...Option) *Service {
	s := &Service{
		l:                l,
		instagramService: instagramService,
		whitelistService: whitelistService,
	}

	allOptions := append([]Option{
		WithCooldown(time.Duration(10 * time.Second)),
	}, options...)

	for _, option := range allOptions {
		option(s)
	}

	return s
}

// Clean ...
func (s *Service) Clean() error {
	users := s.instagramService.GetFollowingUsers()
	var unfollows uint32

	for users.Next() {
		for _, user := range users.Users {
			if s.maxUnfollows != 0 && unfollows == s.maxUnfollows {
				return nil
			}

			isIn, err := s.whitelistService.IsIn(user.ID)
			if err != nil {
				return errors.Wrap(err, "(WhitelistService).IsIn")
			}

			if isIn {
				s.l.Infof("ignoring username %s id %d in whitelist", user.Username, user.ID)
				continue
			}

			if err := user.Unfollow(); err != nil {
				return errors.Wrapf(err, "not able to unfollow user %s id %d", user.Username, user.ID)
			}

			s.l.Infof("unfollowed username %s id %d", user.Username, user.ID)

			unfollows++
			time.Sleep(s.unfollowCooldown)
		}
	}

	s.l.Infof("job unfollowed %d user", unfollows)

	return nil
}
