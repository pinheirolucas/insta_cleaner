package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"github.com/pinheirolucas/insta_cleaner/cleaner"
	"github.com/pinheirolucas/insta_cleaner/helper"
	"github.com/pinheirolucas/insta_cleaner/logger"
	"github.com/pinheirolucas/insta_cleaner/whitelist"
)

func main() {
	var config string
	var limit int

	flag.IntVar(&limit, "limit", 10, "limit of unfollows")
	flag.StringVar(&config, "config", "", "path to config file")
	flag.Parse()

	if limit < 0 {
		log.Fatal("limit must bem greater than 0")
	}

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
	databaseURL := viper.GetString("realtime_database_url")

	session := filepath.Join(sessionsDir, "."+username)

	insta, err := helper.InitLocalGoinsta(username, password, session)
	if err != nil {
		log.Fatalf("helper.InitLocalGoinsta: %v \n", err)
	}

	app, err := firebase.NewApp(
		context.Background(),
		&firebase.Config{
			DatabaseURL: databaseURL,
		},
		option.WithCredentialsFile(firebaseKeyFile),
	)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v \n", err)
	}

	whitelistService, err := whitelist.NewService(context.Background(), viper.GetString("whitelist_service_type"), app)
	if err != nil {
		log.Fatalf("whitelist.NewService: %v \n", err)
	}

	instagramService := cleaner.NewGoinstaInstagramService(insta)
	l := logger.NewDefault()
	service := cleaner.NewService(instagramService, whitelistService, l, cleaner.WithMaxUnfollows(uint32(limit)))

	if err := service.Clean(); err != nil {
		log.Fatalf("(*cleaner.Service).Clean: %v \n", err)
	}
}
