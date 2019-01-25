package cmd

import (
	"fmt"
	"log"

	"github.com/pinheirolucas/insta_cleaner/helper"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get instagram user info bu name",
	Run:   runInfo,
}

func init() {
	userCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		createError("no username provided")
	}

	uname := args[0]

	c, err := loadConfig()
	if err != nil {
		handleError(err)
	}

	insta, err := helper.InitLocalGoinsta(c.username, c.password, c.session)
	if err != nil {
		log.Fatalf("helper.InitLocalGoinsta: %v \n", err)
	}

	user, err := insta.Profiles.ByName(uname)
	if err != nil {
		decorateError(err, "not able to unfollow %s ignoring", uname)
	}

	fmt.Println("id:", user.ID)
	fmt.Println("username:", user.Username)
	fmt.Println("name:", user.FullName)
}
