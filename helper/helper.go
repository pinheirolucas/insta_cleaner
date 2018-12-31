package helper

import (
	"os"

	"github.com/ahmdrz/goinsta"
	"github.com/pkg/errors"
)

// InitLocalGoinsta ...
func InitLocalGoinsta(username, password, session string) (*goinsta.Instagram, error) {
	var insta *goinsta.Instagram

	if _, err := os.Stat(session); err == nil {
		insta, err = goinsta.Import(session)
		if err != nil {
			return nil, errors.Wrap(err, "goinsta.Import")
		}
	} else {
		insta = goinsta.New(username, password)

		if err := insta.Login(); err != nil {
			return nil, errors.Wrap(err, "(*goinsta.Instagram).Login:")
		}

		if err := insta.Export(session); err != nil {
			return nil, errors.Wrap(err, "(*goinsta.Instagram).Export:")
		}
	}

	return insta, nil
}
