package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// This is the default value for version. It is override by ldflags in build time.
var version = ""

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use: "version",
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("v%s\n", version)
}
