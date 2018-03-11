package main

type OauthClient struct {
	ClientSecret string `json:"clientSecret"`
	ClientID     string `json:"clientID"`
	Secret       string `json:"secret"`
}
