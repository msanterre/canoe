package main

import (
	"fmt"
	"net/http"
	"math/rand"

	"github.com/gorilla/mux"
)

func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	slug := generateRandom(8) //TODO: ensure no collision

	redis := RedisClient()
	redis.Cmd("HSET", UrlStore, slug, url)
	fmt.Fprintf(w, fullRedirectionUrl(slug))
}

func ApiRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/submit", SubmitHandler)
	return router
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateRandom(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fullRedirectionUrl(slug string) string {
	return "http://" + redirectionHost() + "/" + slug
}
