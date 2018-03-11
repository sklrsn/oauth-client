package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Index(w http.ResponseWriter, r *http.Request) {
	templates["index.html"].ExecuteTemplate(w, "base", map[string]interface{}{})
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	session, _ := store.Get(r, "session")
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, 302)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Fatalf("Failed to get session object")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		log.Fatalf("Invalid state")
		return
	}

	token, err := oauthConfig.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		log.Fatalf("Incorrect code")
		return
	}

	if !token.Valid() {
		log.Fatalf("Invalid token")
		return
	}

	githubClient := github.NewClient(oauthConfig.Client(oauth2.NoContext, token))
	user, _, err := githubClient.Users.Get(oauth2.NoContext, "")
	if err != nil {
		log.Fatalf("Retrieve userinfo failed")
		return
	}

	session.Values["name"] = user.Name
	session.Values["accessToken"] = token.AccessToken
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
