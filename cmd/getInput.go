package cmd

import (
	"io"
	"os"
)

func getInput(value string) (string, error) {
	if value == "-" {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else if value != "" {
		return value, nil
	}
	return "", nil
}