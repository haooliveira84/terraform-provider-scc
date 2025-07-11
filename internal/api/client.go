package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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

func NewRestApiClient(client *http.Client, baseURL *url.URL, username, password string, caCertBytes []byte, clientCertBytes []byte, clientCertKey []byte) (*RestApiClient, error) {
	useBasicAuth, useCertAuth := isBasicAuthProvided(username, password), isCertAuthProvided(clientCertBytes, clientCertKey)

	if err := validateAuthMode(useBasicAuth, useCertAuth); err != nil {
		return nil, err
	}

	tlsConfig, err := buildTLSConfig(caCertBytes, clientCertBytes, clientCertKey, useCertAuth)
	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		client = withTLSClient(client, tlsConfig)
	} else if client == nil {
		client = &http.Client{}
	}

	return &RestApiClient{
		BaseURL:  baseURL,
		Client:   client,
		Username: username,
		Password: password,
	}, nil
}

func isBasicAuthProvided(username, password string) bool {
	return username != "" && password != ""
}

func isCertAuthProvided(cert, key []byte) bool {
	return len(cert) > 0 && len(key) > 0
}

func validateAuthMode(useBasicAuth, useCertAuth bool) error {
	switch {
	case useBasicAuth && useCertAuth:
		return fmt.Errorf("cannot use both certificate-based and basic authentication simultaneously")
	case !useBasicAuth && !useCertAuth:
		return fmt.Errorf("either certificate-based or basic authentication must be provided")
	default:
		return nil
	}
}

func buildTLSConfig(caCert, clientCert, clientKey []byte, useCertAuth bool) (*tls.Config, error) {
	// if Certificate based authentication, it is mandatory to provide CA Certificate
	if len(caCert) == 0 && !useCertAuth {
		return nil, nil
	}

	tlsConfig := &tls.Config{}

	if len(caCert) > 0 {
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("failed to parse CA certificate: input is not valid PEM-encoded data")
		}
		tlsConfig.RootCAs = caCertPool
	}

	if useCertAuth {
		if err := validatePEMBlock(clientCert, "client certificate"); err != nil {
			return nil, err
		}
		if err := validatePEMBlock(clientKey, "client key"); err != nil {
			return nil, err
		}
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate/key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

func validatePEMBlock(data []byte, label string) error {
	if block, _ := pem.Decode(data); block == nil {
		return fmt.Errorf("%s is not valid PEM-encoded data", label)
	}
	return nil
}

func withTLSClient(client *http.Client, tlsConfig *tls.Config) *http.Client {
	if client == nil {
		client = &http.Client{}
	}
	client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	return client
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

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
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

	// Handle 401 Unauthorized explicitly
	if response.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication rejected: HTTP %d for %s %s. Response: %s",
			response.StatusCode, response.Request.Method, response.Request.URL, string(bodyBytes))
	}

	// Attempt to decode a structured error message
	var errorResp ErrorResponse
	if err := json.Unmarshal(bodyBytes, &errorResp); err == nil && errorResp.Message != "" {
		return fmt.Errorf("HTTP %s %s failed with status %d: %s",
			response.Request.Method, response.Request.URL, response.StatusCode, errorResp.Message)
	}

	// Fallback to raw body
	return fmt.Errorf("HTTP %s %s failed with status %d. Raw response: %s",
		response.Request.Method, response.Request.URL, response.StatusCode, string(bodyBytes))
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
