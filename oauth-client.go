package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const (
	defaultLayout = "templates/layout.html"
	templateDir   = "templates/"

	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	redirectURL        = ""
)

var (
	oauthClient *OauthClient
	oauthConfig *oauth2.Config
	store       *sessions.CookieStore
	scopes      = []string{"repo"}
	templates   = map[string]*template.Template{}
)

func main() {

	oauthClient := OauthClient{ClientID: "xxxxxxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxx", Secret: ""}

	templates["index.html"] = template.Must(template.ParseFiles(
		templateDir+"index.html", defaultLayout))

	store = sessions.NewCookieStore([]byte(oauthClient.Secret))
	oauthConfig = &oauth2.Config{
		ClientID:     oauthClient.ClientID,
		ClientSecret: oauthClient.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeURL,
			TokenURL: githubTokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}

	r := InitializeRouter()
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8088"
	}

	fmt.Println("Server started")
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
