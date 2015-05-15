package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/fzzy/radix/redis"
)

type SubmitResponse struct {
	Url string `json:"url"`
}

func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	redis := RedisClient()
	slug := generateSlug(redis)

	redis.Cmd("HSET", UrlStore, slug, url)

	redirectUrl := fullRedirectionUrl(slug)
	response := &SubmitResponse{Url: redirectUrl}
	submitJsonResponse(w, response)
}

func ApiRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/submit", SubmitHandler)
	return router
}

// Internal methods

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

func slugExists(redis *redis.Client, slug string) bool {
	exists, err := redis.Cmd("HEXISTS", UrlStore, slug).Bool()
	if err != nil {
		fmt.Println("[!ERR] slugExists - ", err)
		// TODO: Better to return true, but it could make infinite loop. Make better.
		return false
	}
	return exists
}

func generateSlug(redis *redis.Client) string {
	slug := generateRandom(8)
	for slugExists(redis, slug) {
		slug = generateRandom(8)
	}
	return slug
}

func submitJsonResponse(w http.ResponseWriter, resp *SubmitResponse) {
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, jsonErr := ToJson(resp)
	if jsonErr != nil {
		http.Error(w, "{\"error\": \"Uknown Error\"}", 500)
		return
	}
	w.Write([]byte(jsonResponse))
}
