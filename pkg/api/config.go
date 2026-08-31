package api

import (
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
)

func Initialize(apiConfig config.API) (string, bool) {
	token = apiConfig.Token
	baseURL = apiURL(apiConfig.Environment)
	if token == "" {
		token = console.Input("Enter your token")
		return token, true
	}
	return token, false
}

func apiURL(environment string) string {
	switch environment {
	case "sandbox":
		return "https://api.sandbox.artifactsmmo.com"
	case "beta":
		return "https://api.beta.artifactsmmo.com"
	default:
		return "https://api.artifactsmmo.com"
	}
}
