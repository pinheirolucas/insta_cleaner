package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/pinheirolucas/insta_cleaner/helper"
)

func main() {
	var config, whitelist string

	flag.StringVar(&config, "config", "", "path to config file")
	flag.StringVar(&whitelist, "whitelist", "", "path to whitelist txt file")
	flag.Parse()

	if whitelist == "" {
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
	session := filepath.Join(sessionsDir, "."+username)

	insta, err := helper.InitLocalGoinsta(username, password, session)
	if err != nil {
		log.Fatalf("helper.InitLocalGoinsta: %v \n", err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(credentials))
	if err != nil {
		log.Fatalf("firebase.NewApp: %v \n", err)
	}

	firestore, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("(*firebase.App).Firestore: %v \n", err)
	}

	file, err := os.Open(whitelist)
	if err != nil {
		log.Fatalf("os.Open: %v \n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		u := scanner.Text()

		user, err := insta.Profiles.ByName(u)
		if err != nil {
			log.Fatalf("(*goinsta.Instagram).Profiles.ByName: %v \n", err)
		}

		snap, err := firestore.Collection("whitelist").Doc(strconv.Itoa(int(user.ID))).Get(context.Background())
		switch grpc.Code(err) {
		case codes.OK:
			if snap.Exists() {
				continue
			}
		case codes.NotFound:
			// proceed do create
		default:
			log.Fatalf("not able to get %s:%d", user.Username, user.ID)
		}

		_, err = snap.Ref.Create(context.Background(), map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		})
		if err != nil {
			log.Fatalf("snap.Ref.Create: %v \n", err)
		}
	}
}
