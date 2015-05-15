package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/fzzy/radix/redis"
	"github.com/asaskevich/govalidator"
)

type SubmitResponse struct {
	Url string `json:"url"`
}

func AuthenticateRequest(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	if requireAuth() {
		authKey := getAuthKey()
		reqAuthKey, _, _ := req.BasicAuth()
		fmt.Println("[auth] (", authKey, "-", reqAuthKey, ")")

		if authKey != reqAuthKey {
			fmt.Println("Not authenticated")
			http.Error(rw, "Authkey invalid", 401)
			return
		}
	}
	next(rw, req)
}

func getAuthKey() string {
  return os.Getenv("AUTH_KEY")
}

func requireAuth() bool {
	if authkey := os.Getenv("AUTH_KEY"); len(authkey) > 1 {
		return true
	}
	return false
}

func SubmitHandler(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	fmt.Println("[submit] - ", url)

	if submitUrlValidation(w, url) {
		redis := RedisClient()
		slug := generateSlug(redis)
		defer redis.Close()

		redis.Cmd("HSET", UrlStore, slug, url)

		redirectUrl := fullRedirectionUrl(slug)
		response := &SubmitResponse{Url: redirectUrl}
		submitJsonResponse(w, response)
	}
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
		http.Error(w, "{\"error\": \"Unknown Error\"}", 500)
		return
	}
	w.Write([]byte(jsonResponse))
}

func submitUrlValidation(w http.ResponseWriter, url string) bool {
	if govalidator.IsURL(url) {
		return true
	}
	http.Error(w, "{\"error\": \"This URL is invalid\"}", 422)
	return false
}
