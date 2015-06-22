package main

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/thoas/stats"
)

const DefaultRedirectionHost = "canoe.lvh.me"
const DefaultAPIHost = "canoeapi.lvh.me"
const DefaultPort = "5050"

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Canoe URL shorterning service")
}

func main() {
	statsMiddleware := stats.New()

	router := mux.NewRouter()
	redirectionRouter := RedirectionRouter()
	apiRouter := ApiRouter(statsMiddleware)

	n := negroni.Classic()

	router.Host(redirectionHost()).Handler(negroni.New(
		statsMiddleware,
		negroni.Wrap(redirectionRouter),
	))

	router.Host(apiHost()).Handler(negroni.New(
		negroni.HandlerFunc(AuthenticateRequest),
		statsMiddleware,
		negroni.Wrap(apiRouter),
	))

	router.HandleFunc("/", HomeHandler)
	n.UseHandler(router)
	n.Run("0.0.0.0:" + runPort())
}

func runPort() string {
	if envPort := os.Getenv("PORT"); len(envPort) > 1 {
		return envPort
	}
	return DefaultPort
}

func apiHost() string {
	if envHost := os.Getenv("CANOE_API_HOST"); len(envHost) > 1 {
		return envHost
	}
	return DefaultAPIHost
}

func redirectionHost() string {
	if envHost := os.Getenv("CANOE_REDIRECTION_HOST"); len(envHost) > 1 {
		return envHost
	}
	return DefaultRedirectionHost
}

func ToJson(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("[ERR] ToJson - ", err)
		return "{}", err
	}
	return string(bytes), nil
}
