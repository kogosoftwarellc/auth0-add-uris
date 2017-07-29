package main

import (
	"github.com/kogosoftwarellc/auth0-add-uris/types"
	"github.com/kogosoftwarellc/go-simple-logging"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/kogosoftwarellc/auth0-add-uris/helpers"
	"os"
)

func main() {
	logger := logging.NewLogger()
	appClientId := kingpin.Flag("app-client-id", "The client id for the app client").
		Required().String()
	clientId := kingpin.Flag("client-id", "The client id for the non interactive client").
		Required().String()
	clientSecret := kingpin.Flag("client-secret", "The client secret for the non interactive client").
		Required().String()
	domain := kingpin.Flag("domain", "The domain of the auth0 app").
		Required().String()
	callbacks := kingpin.Flag("callback", "A callback url to add to the client's allowed callback urls").
		Required().URLList()
	logouts := kingpin.Flag("logout", "A logout url to add to the client's allowed logout urls").
		Required().URLList()
	maxLen := kingpin.Flag("max-len", "The maximum number of urls to allow").
		Default("1").
		Int()

	kingpin.Parse()

	if len(*callbacks) == 0 {
		logger.Error("At least one --callback must be suplied.")
		os.Exit(1)
	}

	if len(*logouts) == 0 {
		logger.Error("At least one --logout must be suplied.")
		os.Exit(1)
	}

	if *maxLen < len(*callbacks) || *maxLen < len(*logouts) {
		logger.Error("maxLen cannot be less than the number of callbacks or logouts")
		os.Exit(1)
	}

	config := types.ApplicationConfig {
		Logger: logger,
		ClientSecret: *clientSecret,
		ClientId: *clientId,
		Callbacks: *callbacks,
		Logouts: *logouts,
		Domain: *domain,
		AppClientId: *appClientId,
		MaxLen: *maxLen,
	}

	helpers.AddUris(config)
}
