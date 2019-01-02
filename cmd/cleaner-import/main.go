package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"github.com/pinheirolucas/insta_cleaner/helper"
	"github.com/pinheirolucas/insta_cleaner/whitelist"
)

func main() {
	var config, wl string

	flag.StringVar(&config, "config", "", "path to config file")
	flag.StringVar(&wl, "whitelist", "", "path to whitelist txt file")
	flag.Parse()

	if wl == "" {
		log.Fatal("no whitelist file")
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
	credentials := viper.GetString("firebase_admin_key_file")
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
		option.WithCredentialsFile(credentials),
	)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v \n", err)
	}

	whitelistService, err := whitelist.NewService(context.Background(), viper.GetString("whitelist_service_type"), app)
	if err != nil {
		log.Fatalf("whitelist.NewService: %v \n", err)
	}

	file, err := os.Open(wl)
	if err != nil {
		log.Fatalf("os.Open: %v \n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		u := scanner.Text()

		log.Printf("importing %s \n", u)

		user, err := insta.Profiles.ByName(u)
		if err != nil {
			log.Fatalf("(*goinsta.Instagram).Profiles.ByName: %v \n", err)
		}

		if err := whitelistService.CreateIfNotExists(user.ID, user.Username); err != nil {
			log.Fatalf("(cleaner.WhitelistService).CreateIfNotExists: %v \n", err)
		}
	}
}
