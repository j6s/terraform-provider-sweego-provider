package sweego

import (
	"net/http"
)

const DefaultBaseUrl = "https://api.sweego.io/"

type SweegoApi struct {
	baseUrl    string
	apiKey     string
	clientId   string
	httpClient *http.Client
	logger     SweegoApiLogger
}

func (api *SweegoApi) WithLogger(logger SweegoApiLogger) *SweegoApi {
	return &SweegoApi{
		baseUrl:    api.baseUrl,
		apiKey:     api.apiKey,
		clientId:   api.clientId,
		httpClient: api.httpClient,
		logger:     logger,
	}
}

func NewSweegoApi(apiKey string, clientId string) SweegoApi {
	return SweegoApi{
		baseUrl:    DefaultBaseUrl,
		apiKey:     apiKey,
		clientId:   clientId,
		httpClient: &http.Client{},
		logger:     GolangLogger{},
	}
}

func NewSweegoApiWithBaseUrl(baseUrl string, apiKey string, clientId string) *SweegoApi {
	return &SweegoApi{
		baseUrl:    baseUrl,
		apiKey:     apiKey,
		clientId:   clientId,
		httpClient: &http.Client{},
		logger:     GolangLogger{},
	}
}
