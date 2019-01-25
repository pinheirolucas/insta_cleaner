package cmd

import (
	"github.com/spf13/cobra"
)

// whitelistCmd represents the whitelist command
var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Instagram users whitelist manager",
}

func init() {
	rootCmd.AddCommand(whitelistCmd)
}
