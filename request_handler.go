package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

	session, _ := store.Get(r, "sess")
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, 302)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sess")
	if err != nil {
		fmt.Fprintln(w, "aborted")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
		return
	}

	tkn, err := oauthConfig.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		fmt.Fprintln(w, "there was an issue getting your token")
		return
	}

	if !tkn.Valid() {
		fmt.Fprintln(w, "retrieved invalid token")
		return
	}

	githubClient := github.NewClient(oauthConfig.Client(oauth2.NoContext, tkn))
	user, _, err := githubClient.Users.Get(oauth2.NoContext, "")
	if err != nil {
		fmt.Println(w, "error getting name")
		return
	}

	session.Values["name"] = user.Name
	session.Values["accessToken"] = tkn.AccessToken
	session.Save(r, w)
	fmt.Println(tkn.AccessToken)
	http.Redirect(w, r, "/", 302)
}
