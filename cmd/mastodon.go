/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mastodonCmd represents the mastodon command
var mastodonCmd = &cobra.Command{
	Use:   "mastodon",
	Short: "Post to Mastodon",
	Long: `Mastodon is LATER`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := GetAlways("MASTODON_SERVER")
		if (err != nil) {
			return err
		}
		fmt.Printf("mastodon posting to %s with %s\n", server, args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mastodonCmd)

	viper.BindEnv("MASTODON_SERVER")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mastodonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mastodonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
