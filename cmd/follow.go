package cmd

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"github.com/pinheirolucas/insta_cleaner/helper"
	"github.com/pinheirolucas/insta_cleaner/whitelist"
)

// followCmd represents the follow command
var followCmd = &cobra.Command{
	Use:   "follow",
	Short: "Follow a list of users by name",
	Run:   runFollow,
}

func init() {
	userCmd.AddCommand(followCmd)
}

func runFollow(cmd *cobra.Command, usernames []string) {
	if len(usernames) == 0 {
		createError("no usernames provided")
	}

	c, err := loadConfig()
	if err != nil {
		handleError(err)
	}

	insta, err := helper.InitLocalGoinsta(c.username, c.password, c.session)
	if err != nil {
		log.Fatalf("helper.InitLocalGoinsta: %v \n", err)
	}

	app, err := firebase.NewApp(
		context.Background(),
		&firebase.Config{
			DatabaseURL: c.databaseURL,
		},
		option.WithCredentialsFile(c.credentials),
	)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v \n", err)
	}

	whitelistService, err := whitelist.NewService(context.Background(), viper.GetString("whitelist_service_type"), app)
	if err != nil {
		log.Fatalf("whitelist.NewService: %v \n", err)
	}

	for _, uname := range usernames {
		log.Printf("following %s \n", uname)

		user, err := insta.Profiles.ByName(uname)
		if err != nil {
			fmt.Printf("not able to follow %s ignoring", uname)
			continue
		}

		if err = user.Follow(); err != nil {
			handleError(err)
		}

		if err := whitelistService.CreateIfNotExists(user.ID, user.Username); err != nil {
			handleError(err)
		}
	}
}
