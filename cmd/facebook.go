package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	facebook_page_access_token string
	facebook_page_id           string
)

var facebookCmd = &cobra.Command{
	Use:   "facebook",
	Short: "Post on a Facebook page",
	Long:  `Post on a Facebook page.  Note: \n * You need to set the FACEBOOK_PAGE_ACCESS_TOKEN in the environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  facebookPost,
}

func InitFacebook() {
	rootCmd.AddCommand(facebookCmd)

	emailCmd.Flags().StringVar(&facebook_page_id, "facebook-page-id", "", "Facebook page's ID")
	viper.BindPFlag("smtp-host", emailCmd.PersistentFlags().Lookup("facebook-page-id"))

	viper.BindEnv("FACEBOOK_PAGE_ACCESS_TOKEN")
}

type facebookPagePost struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

func facebookPost(cmd *cobra.Command, args []string) error {

	facebook_page_id = viper.GetString("facebook_page_id")
	//LATER: are facebook page ID's always numeric?  if so, validate
	facebook_page_access_token = viper.GetString("facebook_page_access_token")

	body, bodyErr := getInput(args[0])
	if bodyErr != nil {
		return bodyErr
	}

	jsonBody := facebookPagePost{
		Message:     body,
		AccessToken: facebook_page_access_token,
	}

	strBody, jsonErr := json.Marshal(jsonBody)
	if jsonErr != nil {
		return jsonErr
	}

	bytesBody := []byte(strBody)
	bodyReader := bytes.NewReader(bytesBody)

	requestURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/feed", facebook_page_id)
	req, reqErr := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Content-Type", "application/json")
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
	slog.Info("Facebook response", "status", res.Status, "body", resBody)
	return nil
}
