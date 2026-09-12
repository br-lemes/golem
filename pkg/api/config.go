package api

import (
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
)

func Initialize(apiConfig config.API) (string, bool) {
	defaultClient = NewClient(apiConfig)
	if defaultClient.Token == "" {
		defaultClient.Token = console.Input("Enter your token")
		return defaultClient.Token, true
	}
	return defaultClient.Token, false
}
