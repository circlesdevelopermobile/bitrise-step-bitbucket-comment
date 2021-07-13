package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func fail(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func checkNon200(response *http.Response) error {
	if response.StatusCode >= 400 {
		return errors.New("Bitbucket status: " + response.Status)
	} else {
		return nil
	}
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
	err = checkNon200(resp)
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

func postComment(client *http.Client, accessToken *AccessToken, config *BitbucketConfig, userData *UserData) error {
	endpt := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/pullrequests/%s/comments", config.RepoSlug, userData.PullRequestId)
	msg, err := userData.getMessage()
	if err != nil {
		return err
	}
	request := CommentRequest{Comment: Comment{Message: *msg}}
	b, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", endpt, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken.BearerToken)
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		return err
	}
	err = checkNon200(resp)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cfg := &BitbucketConfig{}
	usrData := &UserData{}
	httpClient := http.Client{}
	if err := stepconf.Parse(cfg); err != nil {
		fail("Error parson config: %s\n", err)
	}
	if err := stepconf.Parse(usrData); err != nil {
		fail("Error parson config: %s\n", err)
	}
	if len(usrData.PullRequestId) == 0 {
		log.Printf("This is not a PR, skipping")
		os.Exit(0)
	}
	accessToken, err := obtainAccessToken(&httpClient, cfg)
	if err != nil {
		fail("Unable to authenticate to bitbucket: %s\n", err)
	}
	err = postComment(&httpClient, accessToken, cfg, usrData)
	if err != nil {
		fail("Unable to post comment: %s\n", err)
	}
	// Success
	os.Exit(0)
}
