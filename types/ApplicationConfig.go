package types

import (
	"github.com/kogosoftwarellc/go-simple-logging"
	"net/url"
)

type ApplicationConfig struct {
	Logger logging.Logger
	ClientId string
	ClientSecret string
	Callbacks []*url.URL
	Logouts []*url.URL
	Domain string
	AppClientId string
	MaxLen int
}
