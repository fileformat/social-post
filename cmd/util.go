package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetAlways(name string) (string, error) {

	retVal := viper.GetString(name)

	if (retVal == "") {
		return "", fmt.Errorf("%s is not set", name)
	}

	return retVal, nil
}