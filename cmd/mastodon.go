/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	mastodon_server string
)

// mastodonCmd represents the mastodon command
var mastodonCmd = &cobra.Command{
	Use:   "mastodon",
	Short: "Post to Mastodon",
	Long:  `Mastodon is LATER`,
	Args:  cobra.ExactArgs(1),
	RunE:  mastodonPost,
}

func InitMastodon() {
	rootCmd.AddCommand(mastodonCmd)

	mastodonCmd.Flags().StringVar(&mastodon_server, "mastodon-server", "mastodon.social", "Mastodon Server (bare hostname)")
	viper.BindPFlag("mastodon-server", mastodonCmd.PersistentFlags().Lookup("mastodon-server"))

	addImageFlag(mastodonCmd)

	viper.BindEnv("MASTODON_USER_TOKEN")
}

func mastodonPost(cmd *cobra.Command, args []string) error {

	mastodon_server = viper.GetString("mastodon_server")
	if mastodon_server == "" {
		return fmt.Errorf("mastodon server not set")
	}

	mastodon_user_token := viper.GetString("mastodon_user_token")

	hasImage, err := checkImage()
	if err != nil {
		return err
	}

	form := url.Values{}
	if hasImage {
		photoId, photoErr := mastodonPostImage()
		if photoErr != nil {
			return photoErr
		}
		form.Set("media_ids[]", photoId)
		slog.Debug("Mastodon photo ID", "id", photoId)
	}

	body, bodyErr := getInput(args[0])
	if bodyErr != nil {
		return bodyErr
	}

	form.Set("status", body)
	form.Set("visibility", "public")

	requestURL := fmt.Sprintf("https://%s/api/v1/statuses", mastodon_server)
	req, reqErr := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(form.Encode()))
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mastodon_user_token))
	req.Header.Set("User-Agent", UserAgent)
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		return resErr
	}

	resBody, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	slog.Info("Mastodon response", "status", res.Status, "body", resBody)
	return nil
}

type mastodonMedia struct {
	Id string `json:"id"`
}

func mastodonPostImage() (string, error) {

	requestURL := fmt.Sprintf("https://%s/api/v2/media", mastodon_server)

	res, err := postImage(requestURL,
		"file",
		imageFile,
		map[string]string{
			"description": imageCaption,
		},
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", viper.GetString("mastodon_user_token")),
		})
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	slog.Debug("Mastodon media post response", "status", res.Status, "body", string(body))

	var photo mastodonMedia
	jsonErr := json.Unmarshal(body, &photo)
	if jsonErr != nil {
		return "", jsonErr
	}

	return photo.Id, nil
}
