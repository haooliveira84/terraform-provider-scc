package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RestApiClient struct {
	Client   *http.Client
	BaseURL  *url.URL
	Username string
	Password string
}

type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewRestApiClient(client *http.Client, baseURL *url.URL, username, password string, caCertPEM []byte) (*RestApiClient, error) {
	// Create a CA certificate pool and append the CA cert
	var tlsConfig *tls.Config
	if caCertPEM != nil {
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCertPEM); !ok {
			return nil, fmt.Errorf("failed to parse CA certificate: input is not valid PEM-encoded data")
		}

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	// Set up HTTP client and transport
	if client == nil {
		client = &http.Client{}
	}

	if tlsConfig != nil {
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return &RestApiClient{
		BaseURL:  baseURL,
		Client:   client,
		Username: username,
		Password: password,
	}, nil
}

func (c *RestApiClient) DoRequest(method string, endpoint string, body []byte) (*http.Response, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	finalURL := c.BaseURL.ResolveReference(endpointURL)

	req, err := http.NewRequest(method, finalURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()
	err = validateResponse(resp)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func validateResponse(response *http.Response) error {
	if response.StatusCode == http.StatusOK ||
		response.StatusCode == http.StatusCreated ||
		response.StatusCode == http.StatusNoContent {
		return nil
	}

	bodyBytes, _ := io.ReadAll(response.Body)
	var errorResp ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errorResp); err != nil || errorResp.Message == "" {
		return fmt.Errorf("HTTP %s %s failed with status %d. Raw response: %s",
			response.Request.Method, response.Request.URL, response.StatusCode, string(bodyBytes))
	}

	return fmt.Errorf("HTTP %s %s failed with status %d: %s",
		response.Request.Method, response.Request.URL, response.StatusCode, errorResp.Message)
}

func (c *RestApiClient) GetRequest(endpoint string) (*http.Response, error) {
	return c.DoRequest(http.MethodGet, endpoint, nil)
}

func (c *RestApiClient) PostRequest(endpoint string, body []byte) (*http.Response, error) {
	return c.DoRequest(http.MethodPost, endpoint, body)
}

func (c *RestApiClient) PutRequest(endpoint string, body []byte) (*http.Response, error) {
	return c.DoRequest(http.MethodPut, endpoint, body)
}

func (c *RestApiClient) DeleteRequest(endpoint string) (*http.Response, error) {
	return c.DoRequest(http.MethodDelete, endpoint, nil)
}
