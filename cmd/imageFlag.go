package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	imageFile    string
	imageCaption string
)

func addImageFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&imageFile, "image", "", "Image file")
	viper.BindPFlag("image", cmd.PersistentFlags().Lookup("image"))

	cmd.Flags().StringVar(&imageCaption, "image-caption", "", "Image file")
	viper.BindPFlag("image-caption", cmd.PersistentFlags().Lookup("image-caption"))

}

func checkImage() (bool, error) {

	// check if the image file is specified
	if imageFile == "" {
		return false, nil
	}

	// check if the image file exists
	if _, err := os.Stat(imageFile); os.IsNotExist(err) {
		return false, fmt.Errorf("image file %s does not exist", imageFile)
	}

	//LATER: check image format

	return true, nil
}

func postImage(dst, fieldname, filename string, otherFields map[string]string, otherHeaders map[string]string) (*http.Response, error) {

	f, _ := os.Open(filename)
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, filepath.Base(f.Name()))
	if err != nil {
		return nil, err
	}

	io.Copy(part, f)

	for key, val := range otherFields {
		err = writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}
	writer.Close()

	req, err := http.NewRequest("POST", dst, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("User-Agent", UserAgent)
	for key, val := range otherHeaders {
		req.Header.Add(key, val)
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	slog.Info("Uploading image", "filename", imageFile)
	return client.Do(req)
}
