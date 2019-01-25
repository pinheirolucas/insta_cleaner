package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	username    string
	password    string
	sessionsDir string
	credentials string
	databaseURL string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "icleaner",
	Short: "Toolbelt for insta_cleaner",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.icleaner.yaml)")

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "instagram username")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "instagram password")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().StringVarP(&sessionsDir, "sessions-dir", "s", "", "session dump target dir")
	viper.BindPFlag("sessionsDir", rootCmd.PersistentFlags().Lookup("session-dir"))

	rootCmd.PersistentFlags().StringVarP(&credentials, "credentials", "c", "", "firebase admin credentials file location")
	viper.BindPFlag("firebase_admin_key_file", rootCmd.PersistentFlags().Lookup("credentials"))

	rootCmd.PersistentFlags().StringVar(&databaseURL, "db", "", "firebase realtime db URL")
	viper.BindPFlag("realtime_database_url", rootCmd.PersistentFlags().Lookup("db"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			handleError(err)
		}

		pwd, err := os.Getwd()
		if err != nil {
			handleError(err)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(pwd)
		viper.SetConfigName(".icleaner")
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func decorateError(err error, msg string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Println(errors.Wrap(err, msg))
	} else {
		fmt.Println(errors.Wrapf(err, msg, args...))
	}

	os.Exit(1)
}

func createError(err string) {
	fmt.Println(err)
	os.Exit(1)
}
