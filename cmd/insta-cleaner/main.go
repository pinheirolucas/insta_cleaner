package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gopkg.in/ahmdrz/goinsta.v2"
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

	users := insta.Account.Following()
	for users.Next() {
		for _, user := range users.Users {
			_, err := firestore.Collection("whitelist").Doc(strconv.Itoa(int(user.ID))).Get(context.Background())
			switch grpc.Code(err) {
			case codes.OK:
				log.Printf("ignoring username %s id %d in whitelist", user.Username, user.ID)
				continue
			case codes.NotFound:
				// continue to unfollow user
			default:
				log.Fatalf("not able to Get doc ID: %d: %v \n", user.ID, err)
			}

			if err := user.Unfollow(); err != nil {
				log.Fatalf("not able to unfollow %s id %d: %v \n", user.Username, user.ID, err)
			}

			log.Printf("user unfollowed: %s", user.Username)

			time.Sleep(time.Duration(3 * time.Second))
		}
	}
}
