
package main

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
)

const DefaultRedirectionHost = "canoe.lvh.me"
const DefaultAPIHost = "canoeapi.lvh.me"
const DefaultPort = "5050"

func main() {
	// TODO: Ensure redis is running & show stats
	router := mux.NewRouter()
	redirectionRouter := RedirectionRouter()
	apiRouter := ApiRouter()

	n := negroni.Classic()

	router.Host(redirectionHost()).Handler(negroni.New(
		negroni.Wrap(redirectionRouter),
	))

	router.Host(apiHost()).Handler(negroni.New(
		negroni.Wrap(apiRouter),
	))

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
