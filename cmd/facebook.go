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

	facebookCmd.Flags().StringVar(&facebook_page_id, "facebook-page-id", "", "Facebook page's ID")
	viper.BindPFlag("facebook-page-id", facebookCmd.PersistentFlags().Lookup("facebook-page-id"))

	addImageFlag(facebookCmd)

	viper.BindEnv("FACEBOOK_PAGE_ACCESS_TOKEN")
}

type facebookAttachedMedia struct {
	Id string `json:"media_fbid"`
}

type facebookPagePost struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	Media []facebookAttachedMedia `json:"attached_media,omitempty"`
}

func facebookPost(cmd *cobra.Command, args []string) error {

	facebook_page_id = viper.GetString("facebook_page_id")
	//LATER: are facebook page ID's always numeric?  if so, validate
	facebook_page_access_token = viper.GetString("facebook_page_access_token")

	hasImage, err := checkImage()
	if err != nil {
		return err
	}

	var media facebookAttachedMedia
	if hasImage {
		photoId, photoErr := facebookPostImage()
		if photoErr != nil {
			return photoErr
		}
		media = facebookAttachedMedia{
			Id: photoId,
		}
		slog.Debug("Facebook photo ID", "id", photoId)
	}

	body, bodyErr := getInput(args[0])
	if bodyErr != nil {
		return bodyErr
	}

	jsonBody := facebookPagePost{
		Message:     body,
		AccessToken: facebook_page_access_token,
		Media: []facebookAttachedMedia{media},
	}

	strBody, jsonErr := json.Marshal(jsonBody)
	if jsonErr != nil {
		return jsonErr
	}
	slog.Debug("Facebook request", "body", string(strBody))

	bytesBody := []byte(strBody)
	bodyReader := bytes.NewReader(bytesBody)

	requestURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/feed", facebook_page_id)
	req, reqErr := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Content-Type", "application/json")
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
	slog.Info("Facebook response", "status", res.Status, "body", resBody)
	return nil
}

type facebookPhoto struct {
	Id string `json:"id"`
}

func facebookPostImage() (string, error) {

	requestURL := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/photos", facebook_page_id)

	res, err := postImage(requestURL,
		"source",
		imageFile,
		map[string]string{
			"published": "false",
			"access_token": facebook_page_access_token,
			"caption": imageCaption,
		})
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var photo facebookPhoto
	jsonErr := json.Unmarshal(body, &photo)
	if jsonErr != nil {
		return "", jsonErr
	}

	return photo.Id, nil
}
