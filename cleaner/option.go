package cleaner

import "time"

// Option ...
type Option func(s *Service)

// WithCooldown ...
func WithCooldown(cooldown time.Duration) Option {
	return func(s *Service) {
		s.unfollowCooldown = cooldown
	}
}

// WithMaxUnfollows ...
func WithMaxUnfollows(max uint32) Option {
	return func(s *Service) {
		s.maxUnfollows = max
	}
}
