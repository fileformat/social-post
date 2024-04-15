package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	slack_channel   string
	slack_bot_token string
)

var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Post to a Slack channel",
	Long:  `Post on a Slack channel.  Note: \n * You need to set the SLACK_BOT_TOKEN in the environment`,
	Args:  cobra.ExactArgs(1),
	RunE:  slackPost,
}

func InitSlack() {
	rootCmd.AddCommand(slackCmd)

	slackCmd.Flags().StringVar(&slack_channel, "channel", "", "Slack Channel ID")

	addImageFlag(slackCmd)

	viper.BindEnv("SLACK_BOT_TOKEN")
}

type slackSectionText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type slackSectionBlock struct {
	Type string           `json:"type"`
	Text slackSectionText `json:"text"`
}

type slackSlackFile struct {
	Id string `json:"id"`
}

type slackImageBlock struct {
	Type    string         `json:"type"`
	Image   slackSlackFile `json:"slack_file"`
	AltText string         `json:"alt_text"`
}

type slackPostMessage struct {
	Channel string        `json:"channel"`
	Text    string        `json:"text,omitempty"`
	Blocks  []interface{} `json:"blocks"`
}

func slackPost(cmd *cobra.Command, args []string) error {

	slack_bot_token = viper.GetString("slack_bot_token")

	hasImage, err := checkImage()
	if err != nil {
		return err
	}

	body, bodyErr := getInput(args[0])
	if bodyErr != nil {
		return bodyErr
	}

	blocks := make([]interface{}, 0)

	sectionBlock := slackSectionBlock{
		Type: "section",
		Text: slackSectionText{
			Type: "mrkdwn",
			Text: body,
		},
	}
	blocks = append(blocks, sectionBlock)

	if hasImage {
		photoId, photoErr := slackUploadImage()
		if photoErr != nil {
			return photoErr
		}
		blocks = append(blocks, slackImageBlock{
			Type: "image",
			Image: slackSlackFile{
				Id: photoId,
			},
			AltText: "Image",
		})
	}

	message := slackPostMessage{
		Channel: slack_channel,
		Text:    fmt.Sprintf("Plain text: %s", body),
		Blocks:  blocks,
	}

	bytesBody, jsonErr := json.Marshal(message)
	if jsonErr != nil {
		return jsonErr
	}
	slog.Info("Slack body", "body", string(bytesBody))

	bodyReader := bytes.NewReader(bytesBody)

	req, reqErr := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bodyReader)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slack_bot_token))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
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
	slog.Info("Slack response", "status", res.Status, "body", resBody)
	return nil
}

type slackGetUploadURLExternal struct {
	Ok        bool   `json:"ok"`
	UploadURL string `json:"upload_url"`
	FileId    string `json:"file_id"`
}

func slackUploadImage() (string, error) {

	imageStat, statErr := os.Stat(imageFile)
	if statErr != nil {
		return "", statErr
	}

	requestURL := fmt.Sprintf("https://slack.com/api/files.getUploadURLExternal?filename=%s&length=%d", imageFile, imageStat.Size())
	urlReq, urlReqErr := http.NewRequest("GET", requestURL, nil)
	if urlReqErr != nil {
		return "", urlReqErr
	}
	urlReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slack_bot_token))
	urlReq.Header.Set("User-Agent", UserAgent)
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	urlResp, urlRespErr := client.Do(urlReq)
	if urlRespErr != nil {
		return "", urlRespErr
	}
	var uploadURL = slackGetUploadURLExternal{}
	uploadURLBody, uploadURLErr := io.ReadAll(urlResp.Body)
	if uploadURLErr != nil {
		return "", uploadURLErr
	}
	jsonErr := json.Unmarshal(uploadURLBody, &uploadURL)
	if jsonErr != nil {
		return "", jsonErr
	}

	if !uploadURL.Ok {
		return "", fmt.Errorf("failed to get upload URL: %s", uploadURLBody)
	}

	f, _ := os.Open(imageFile)
	defer f.Close()

	// post the image to the upload URL
	postReq, postReqErr := http.NewRequest("POST", uploadURL.UploadURL, f)
	if postReqErr != nil {
		return "", postReqErr
	}
	postReq.Header.Set("Content-Type", "image/jpeg")
	postReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slack_bot_token))
	postReq.Header.Set("User-Agent", UserAgent)

	_, postErr := client.Do(postReq)
	if postErr != nil {
		return "", postErr
	}

	// complete the upload
	completeBody := fmt.Sprintf("{\"files\": [{\"id\":\"%s\"}]}", uploadURL.FileId)
	completeReq, completeReqErr := http.NewRequest("POST", "https://slack.com/api/files.completeUploadExternal", strings.NewReader(completeBody))
	if completeReqErr != nil {
		return "", completeReqErr
	}
	completeReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	completeReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slack_bot_token))
	completeReq.Header.Set("User-Agent", UserAgent)

	_, completeRespErr := client.Do(completeReq)
	if completeRespErr != nil {
		return "", completeRespErr
	}

	//discovered after much pain: sleep for 1 second
	time.Sleep(5 * time.Second)

	slog.Info("Uploaded image", "file_id", uploadURL.FileId)
	return uploadURL.FileId, nil
}
