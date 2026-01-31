package sweego

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (api *SweegoApi) executeRequest(
	method string,
	endpoint string,
	body interface{},
	serializationType string,
	responseData interface{},
) error {
	headers := map[string]string{}
	absUrl := fmt.Sprintf("%s/%s", strings.TrimRight(api.baseUrl, "/"), strings.TrimLeft(endpoint, "/"))

	headers["Accept"] = "application/json"
	headers["Api-Key"] = api.apiKey

	var bodyString string
	if serializationType == "json" {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("Error executing request %s %s: Cannot serialize JSON Body: %s", method, absUrl, err)
		}

		bodyString = string(bodyBytes)
		headers["Content-Type"] = "application/json"
	} else if serializationType == "form" {
		data := url.Values{}
		for k, v := range headers {
			data.Set(k, v)
		}
		bodyString = data.Encode()
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else if serializationType == "none" || serializationType == "" {
		bodyString = ""
	} else {
		return fmt.Errorf("Error executing request %s %s: Unknown serialization type: %s", method, absUrl, serializationType)
	}

	request, err := http.NewRequest(method, absUrl, strings.NewReader(bodyString))
	if err != nil {
		return fmt.Errorf("Error executing request %s %s: Cannot initialize request: %s", method, absUrl, err)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	api.logger.Debug(fmt.Sprintf("%s %s\n%#v", request.Method, request.URL, request))
	response, err := api.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Error executing request %s %s: Error executing request: %s", method, absUrl, err)
	}
	api.logger.Debug(fmt.Sprintf("%#v", response))
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error executing request %s %s: Cannot read response body: %s", method, absUrl, err)
	}

	api.logger.Debug(string(responseBody))

	if responseData != nil {

		if response.StatusCode >= 300 {
			return fmt.Errorf("Error executing request %s %s: Invalid response status code: %d\n%s\n", method, absUrl, response.StatusCode, responseBody)
		}

		err = json.Unmarshal(responseBody, responseData)
		if err != nil {
			return fmt.Errorf("Error executing request %s %s: Cannot parse response body: %s\n%s", method, absUrl, err, responseBody)
		}
	}

	return nil
}

func (api *SweegoApi) executeJsonRequest(method string, endpoint string, body interface{}, responseData interface{}) error {
	return api.executeRequest(method, endpoint, body, "json", responseData)
}

func (api *SweegoApi) executePlainRequest(method string, endpoint string, responseData interface{}) error {
	return api.executeRequest(method, endpoint, nil, "", responseData)
}

func (api *SweegoApi) executeGetRequest(endpoint string, responseData interface{}) error {
	return api.executeRequest("GET", endpoint, nil, "", responseData)
}
