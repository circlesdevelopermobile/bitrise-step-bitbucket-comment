package main

import (
	"github.com/bitrise-io/go-steputils/stepconf"
	"os"
	"strings"
)

// BitbucketConfig Client Credentials Grant
// https://developer.atlassian.com/cloud/bitbucket/oauth-2/
type BitbucketConfig struct {
	ClientId stepconf.Secret `env:"bitbucket_client_id,required"`
	Secret   stepconf.Secret `env:"bitbucket_client_secret,required"`
	RepoSlug string          `env:"bitbucket_repo_slug,required"`
}

type UserData struct {
	PullRequestId   string `env:"bitbucket_pr_id,required"`
	Message         string `env:"step_message_text"`
	MessageFilePath string `env:"step_message_file"`
}

// Returns the contents of MessageFilePath if that exists
func (this *UserData) getMessage() (*string, error) {
	if len(this.MessageFilePath) > 0 {
		b, err := os.ReadFile(this.MessageFilePath)
		if err != nil {
			return nil, err
		}
		var builder strings.Builder
		builder.WriteString("```\n")
		builder.Write(b)
		builder.WriteString("\n```")
		out := builder.String()
		return &out, nil
	} else {
		return &this.Message, nil
	}
}

// AccessToken - Bitbucket response format
type AccessToken struct {
	Scopes       string `json:"scopes"`
	BearerToken  string `json:"access_token"`
	Expiry       int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	State        string `json:"client_credentials"`
	RefreshToken string `json:"refresh_token"`
}

// CommentRequest - Bitbucket request format
type CommentRequest struct {
	Comment Comment `json:"content"`
}

// Comment - Nested object for the actual message
type Comment struct {
	Message string `json:"raw"`
}
