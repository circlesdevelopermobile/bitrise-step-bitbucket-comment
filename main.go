package main

import (
	"encoding/json"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// BitbucketConfig Client Credentials Grant
// https://developer.atlassian.com/cloud/bitbucket/oauth-2/
type BitbucketConfig struct {
	BaseUrl  string          `env:"bitbucket_base_url,required"`
	ClientId stepconf.Secret `env:"bitbucket_client_id,required"`
	Secret   stepconf.Secret `env:"bitbucket_client_secret,required"`
}

func fail(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func obtainAccessToken(client *http.Client, config *BitbucketConfig) (*AccessToken, error) {
	endpt := "https://bitbucket.org/site/oauth2/access_token"
	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")

	req, _ := http.NewRequest("POST", endpt, strings.NewReader(payload.Encode()))
	req.SetBasicAuth(string(config.ClientId), string(config.Secret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func main() {
	var cfg = &BitbucketConfig{}
	httpClient := http.Client{}
	if err := stepconf.Parse(cfg); err != nil {
		fail("Error parson config: %s\n", err)
	}
	stepconf.Print(cfg)
	accessToken, err := obtainAccessToken(&httpClient, cfg)
	if err != nil {
		fail("Unable to authenticate to bitbucket: %s\n", err)
	}
}
