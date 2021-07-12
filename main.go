package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	BaseUrl       string          `env:"bitbucket_base_url,required"`
	ClientId      stepconf.Secret `env:"bitbucket_client_id,required"`
	Secret        stepconf.Secret `env:"bitbucket_client_secret,required"`
	RepoSlug      string          `env:"bitbucket_repo_slug,required"`
	PullRequestID string          `env:"BITRISE_PULL_REQUEST"`
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
	body, _ := ioutil.ReadAll(resp.Body)
	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func postComment(client *http.Client, accessToken *AccessToken, config *BitbucketConfig) {
	endpt := fmt.Sprintf("%s/repositories/%s/pullrequests/%s/comments", config.BaseUrl, config.RepoSlug, config.PullRequestID)
	request := CommentRequest{Comment: Comment{Message: "another test"}}
	b, _ := json.Marshal(request)
	println(string(b))

	req, _ := http.NewRequest("POST", endpt, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken.BearerToken)
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Errorf("error: %s\n", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	println(string(body))
}

func main() {
	var cfg = &BitbucketConfig{}
	httpClient := http.Client{}
	//if err := stepconf.Parse(cfg); err != nil {
	//	fail("Error parson config: %s\n", err)
	//}
	stepconf.Print(cfg)
	if len(cfg.PullRequestID) == 0 {
		log.Printf("This is not a PR, skipping")
		os.Exit(0)
	}
	accessToken, err := obtainAccessToken(&httpClient, cfg)
	if err != nil {
		fail("Unable to authenticate to bitbucket: %s\n", err)
	}
	postComment(&httpClient, accessToken, cfg)
}
