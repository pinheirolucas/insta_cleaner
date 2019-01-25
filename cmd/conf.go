package cmd

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type config struct {
	username    string
	password    string
	sessionsDir string
	credentials string
	databaseURL string
	session     string
}

func loadConfig() (*config, error) {
	c := new(config)

	c.username = viper.GetString("user")
	c.password = viper.GetString("password")
	c.sessionsDir = viper.GetString("sessions_dir")
	c.credentials = viper.GetString("firebase_admin_key_file")
	c.databaseURL = viper.GetString("realtime_database_url")
	c.session = filepath.Join(c.sessionsDir, "."+c.username)

	return c, nil
}
