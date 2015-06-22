package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"errors"
	"strings"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/fzzy/radix/redis"
	"github.com/thoas/stats"
)

type SubmitResponse struct {
	Url string `json:"url"`
}

func AuthenticateRequest(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	if requireAuth() {
		authKey := getAuthKey()
		reqAuthKey := requestAuthKey(req)

		fmt.Println("[auth] (", authKey, "-", reqAuthKey, ")")

		if authKey != reqAuthKey {
			fmt.Println("Not authenticated")
			http.Error(rw, "Authkey invalid", 401)
			return
		}
	}
	next(rw, req)
}

func requestAuthKey(req *http.Request) string {
	authKey, _, ok := req.BasicAuth()
	if ok {
		return authKey
	}
	return req.FormValue("apikey")
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
	submitUrl, _ := url.QueryUnescape(req.FormValue("url"))

	fmt.Println("[submit] - ", submitUrl)

	if submitUrlValidation(w, submitUrl) {
		redis := RedisClient()
		defer redis.Close()

		ensureHTTP(&submitUrl)
		slug, err := getSlug(redis, w, req)

		if err == nil {
			fmt.Println("Giving slug: ", slug)
			redis.Cmd("HSET", UrlStore, slug, submitUrl)
			redirectUrl := fullRedirectionUrl(slug)
			response := &SubmitResponse{Url: redirectUrl}
			submitJsonResponse(w, response)
		}
	}
}

func ApiRouter(statsMiddleware *stats.Stats) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/submit", SubmitHandler)
	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := statsMiddleware.Data()
		b, _ := json.Marshal(stats)
		w.Write(b)
	})

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

func getSlug(redis *redis.Client, w http.ResponseWriter, r *http.Request) (string, error) {
	if slugParam := r.FormValue("slug"); len(slugParam) > 1 {
		if slugExists(redis, slugParam) {
			http.Error(w, "{\"error\": \"This slug is already taken.\"}", 422)
			return slugParam, errors.New("Slug taken")
		}
		return slugParam, nil
	}
	return generateSlug(redis), nil
}

func generateSlug(redis *redis.Client) string {
	slug := generateRandom(8)
	for slugExists(redis, slug) {
		slug = generateRandom(8)
	}
	return slug
}

func ensureHTTP(str *string) {
	if strings.Contains(*str, "http") {
		return
	}

	*str = "http://" + *str
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

func submitUrlValidation(w http.ResponseWriter, submitUrl string) bool {
	if strings.Contains(submitUrl, ".") {
		return true
	}
	http.Error(w, "{\"error\": \"This URL is invalid\"}", 422)
	return false
}
