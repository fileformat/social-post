/*
Copyright © 2023 fileformat@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "social-post",
	Short: "Post to social networks",
	Long: `Post an image and text to various social media sites.  See the specific subcommands for details.`,
	//SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	slog.Info("root.Execute")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.social-post.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".social-post" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".social-post")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Debug("Using config file", "file", viper.ConfigFileUsed())
	} else {
		slog.Debug("No config file found", "err", err)	
	}
}
