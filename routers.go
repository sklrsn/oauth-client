package main

import (
	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", Index)
	router.HandleFunc("/redirect", Redirect)
	router.HandleFunc("/callback", Callback)
	return router
}
