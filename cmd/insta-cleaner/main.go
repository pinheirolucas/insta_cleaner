package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/ahmdrz/goinsta"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"github.com/pinheirolucas/insta_cleaner/cleaner"
	"github.com/pinheirolucas/insta_cleaner/logger"
)

func main() {
	var config string

	flag.StringVar(&config, "config", "", "path to config file")

	if config == "" {
		viper.SetConfigName(".insta-cleaner")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("viper.ReadInConfig: %v \n", err)
		}
	} else {
		file, err := os.Open(config)
		if err != nil {
			log.Fatalf("os.Open: %v \n", err)
		}

		if err := viper.ReadConfig(file); err != nil {
			log.Fatalf("viper.ReadConfig: %v \n", err)
		}
	}

	username := viper.GetString("username")
	password := viper.GetString("password")
	sessionsDir := viper.GetString("sessions_dir")
	firebaseKeyFile := viper.GetString("firebase_admin_key_file")
	sessionFile := filepath.Join(sessionsDir, "."+username)

	var insta *goinsta.Instagram

	if _, err := os.Stat(sessionFile); err == nil {
		insta, err = goinsta.Import(sessionFile)
		if err != nil {
			log.Fatalf("goinsta.Import: %v \n", err)
		}
	} else {
		insta = goinsta.New(username, password)

		if err := insta.Login(); err != nil {
			log.Fatalf("(*goinsta.Instagram).Login: %v \n", err)
		}

		if err := insta.Export(sessionFile); err != nil {
			log.Fatalf("(*goinsta.Instagram).Export: %v \n", err)
		}
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(firebaseKeyFile))
	if err != nil {
		log.Fatalf("firebase.NewApp: %v \n", err)
	}

	firestore, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("(*firebase.App).Firestore: %v \n", err)
	}

	instagramService := cleaner.NewGoinstaInstagramService(insta)
	whitelistService := cleaner.NewFirebaseWhitelistService(firestore)
	l := logger.NewDefault()
	service := cleaner.NewService(instagramService, whitelistService, l)

	if err := service.Clean(); err != nil {
		log.Fatalf("(*cleaner.Service).Clean: %v \n", err)
	}
}
