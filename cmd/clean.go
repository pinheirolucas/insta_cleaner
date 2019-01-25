package cmd

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/pinheirolucas/insta_cleaner/cleaner"
	"github.com/pinheirolucas/insta_cleaner/helper"
	"github.com/pinheirolucas/insta_cleaner/logger"
	"github.com/pinheirolucas/insta_cleaner/whitelist"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var limit uint32

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Starts the instagram clean job",
	Run:   runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().Uint32VarP(&limit, "limit", "l", 10, "limit of unfollows")
}

func runClean(cmd *cobra.Command, args []string) {
	if limit < 0 {
		createError("limit must bem greater than 0")
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

	instagramService := cleaner.NewGoinstaInstagramService(insta)
	l := logger.NewDefault()
	service := cleaner.NewService(instagramService, whitelistService, l, cleaner.WithMaxUnfollows(limit))

	if err := service.Clean(); err != nil {
		log.Fatalf("(*cleaner.Service).Clean: %v \n", err)
	}
}
