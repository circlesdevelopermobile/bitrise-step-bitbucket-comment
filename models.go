package main

type AccessToken struct {
	Scopes       string `json:"scopes"`
	BearerToken  string `json:"access_token"`
	Expiry       int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	State        string `json:"client_credentials"`
	RefreshToken string `json:"refresh_token"`
}
type CommentRequest struct {
	Comment Comment `json:"content"`
}

type Comment struct {
	Message string `json:"raw"`
}
