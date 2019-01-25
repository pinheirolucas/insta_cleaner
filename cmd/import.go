package cmd

import (
	"bufio"
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"github.com/pinheirolucas/insta_cleaner/helper"
	"github.com/pinheirolucas/insta_cleaner/whitelist"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Add the users to database from a whitelist",
	Run:   runImport,
}

func init() {
	whitelistCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		createError("no whitelist file")
	}

	wl := args[0]

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
