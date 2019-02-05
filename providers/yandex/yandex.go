package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// @see https://tech.yandex.ru/translate/doc/dg/concepts/about-docpage/

type (
	// Client provider engine
	Client struct {
		apiKey, apiVersion string
		httpClient         *http.Client
	}
)

// New create and return new yandex translate client
func New(key, apiVersion string) *Client {
	return &Client{
		apiKey:     key,
		apiVersion: apiVersion,
		httpClient: &http.Client{},
	}
}

func (c *Client) getBaseURL() string {
	return fmt.Sprintf("%s/%s/%s/translate", baseURL, c.apiVersion, jsonInterface)
}

// Translate return translate for text
func (c *Client) Translate(text, lang string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, c.getBaseURL(), nil)
	if err != nil {
		return nil, err
	}

	u := url.URL{}
	payload := u.Query()
	payload.Add("key", c.apiKey)
	payload.Add("lang", lang)
	payload.Add("text", text)
	payload.Add("format", Plain.String())

	req.URL.RawQuery = payload.Encode()

	raw, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := raw.Body.Close(); err != nil {
			log.Warn("yandex: can't close response body")
		}
	}()

	response := &TranslateResponse{}
	err = json.Unmarshal(b, &response)
	return response.Text, nil
}
