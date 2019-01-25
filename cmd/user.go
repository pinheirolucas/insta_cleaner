package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Instagram user related commands",
}

func init() {
	rootCmd.AddCommand(userCmd)
}
