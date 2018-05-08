package uzan

import (
	"errors"
	"os"
	"testing"
)

var accessToken = os.Getenv("YouZanAccessToken")

var client = ZanClient{IsOAuth: true, AccessToken: accessToken}

func init() {
	if accessToken == "" {
		panic(errors.New("Please check access token"))
	}
}
