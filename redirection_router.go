package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RedirectionHandler(w http.ResponseWriter, req *http.Request) {
	slug := mux.Vars(req)["slug"]
	url := getURLFromSlug(slug)

	if len(url) < 1 {
		http.NotFound(w, req)
		return
	}

	http.Redirect(w, req, url, 301)
}

func RedirectionRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{slug}", RedirectionHandler)

	return router
}

func getURLFromSlug(slug string) string {
  redis := RedisClient()
	url, err := redis.Cmd("HGET", UrlStore, slug).Str()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return url
}
