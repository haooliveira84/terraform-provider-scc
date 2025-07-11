package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestRestApiClient_BasicAuth(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "testuser" || password != "testpassword" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := createBasicAuthClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create basic auth client: %v", err)
	}

	t.Run("GET /success", func(t *testing.T) {
		resp, err := client.GetRequest("/success")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})
}

func TestRestApiClient_CertificateAuth(t *testing.T) {
	// Generate server cert
	serverCertPEM, serverKeyPEM, _, err := generateSelfSignedCert()
	if err != nil {
		t.Fatalf("server cert generation failed: %v", err)
	}

	// Generate client cert
	clientCertPEM, clientKeyPEM, clientCert, err := generateSelfSignedCert()
	if err != nil {
		t.Fatalf("client cert generation failed: %v", err)
	}

	// Create server that validates client cert
	clientCertPool := x509.NewCertPool()
	clientCertPool.AddCert(clientCert)

	serverTLSCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		t.Fatalf("invalid server cert/key: %v", err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/secure", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "secured"}`))
	})

	server := httptest.NewUnstartedServer(handler)
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	server.StartTLS()
	defer server.Close()

	client, err := createCertAuthClient(server.URL, serverCertPEM, clientCertPEM, clientKeyPEM)
	if err != nil {
		t.Fatalf("failed to create cert auth client: %v", err)
	}

	t.Run("GET /secure", func(t *testing.T) {
		resp, err := client.GetRequest("/secure")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})
}

func TestRestApiClient_BothAuthProvidedFails(t *testing.T) {
	baseURL, _ := url.Parse("https://localhost")
	certPEM, keyPEM, _, _ := generateSelfSignedCert()
	// Provided both basic authentication(username/password) and certificate based authentication to the function
	_, err := NewRestApiClient(nil, baseURL, "user", "pass", certPEM, certPEM, keyPEM)
	if err == nil || err.Error() != "cannot use both certificate-based and basic authentication simultaneously" {
		t.Fatalf("expected error for both auth methods provided, got: %v", err)
	}
}

func TestRestApiClient_NoAuthProvidedFails(t *testing.T) {
	baseURL, _ := url.Parse("https://localhost")
	// Provided neither basic authentication(username/password) nor certificate based authentication to the function
	_, err := NewRestApiClient(nil, baseURL, "", "", nil, nil, nil)
	if err == nil || err.Error() != "either certificate-based or basic authentication must be provided" {
		t.Fatalf("expected error for no auth provided, got: %v", err)
	}
}

func TestRestApiClient_InvalidClientCertFails(t *testing.T) {
	baseURL, _ := url.Parse("https://localhost")
	// Generate invalid client certificate and key and provided to the function
	invalidPEM := []byte("not a valid pem")
	_, err := NewRestApiClient(nil, baseURL, "", "", nil, invalidPEM, invalidPEM)
	if err == nil || err.Error() != "client certificate is not valid PEM-encoded data" {
		t.Fatalf("expected PEM validation error, got: %v", err)
	}
}

func TestRestApiClient_InvalidCACertFails(t *testing.T) {
	baseURL, _ := url.Parse("https://localhost")
	// Generate valid client certificate and key
	certPEM, keyPEM, _, _ := generateSelfSignedCert()
	// Generate invalid CA Certificate
	invalidCA := []byte("not valid pem")
	_, err := NewRestApiClient(nil, baseURL, "", "", invalidCA, certPEM, keyPEM)
	if err == nil || err.Error() != "failed to parse CA certificate: input is not valid PEM-encoded data" {
		t.Fatalf("expected CA cert parse error, got: %v", err)
	}
}

// generateSelfSignedCert generates a self-signed TLS certificate and its private key.
func generateSelfSignedCert() (certPEM, keyPEM []byte, cert *x509.Certificate, err error) {
	// Generate a new RSA private key with 2048-bit length
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a certificate template with required fields
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,

		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:    []string{"localhost"},
	}

	// Create a self-signed certificate using the template and the generated private key.
	// The template is being used for parent issuer certificate and the certificate itself since it is a self-signed certfificate.
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Encode the certificate & private key to PEM format
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})

	// Parse the DER-encoded certificate into an x509.Certificate object
	parsedCert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, nil, err
	}

	// Return the PEM-encoded certificate, key, and parsed certificate
	return certPEM, keyPEM, parsedCert, nil
}

func TestValidateResponse_SuccessCodes(t *testing.T) {
	successStatuses := []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}

	for _, code := range successStatuses {
		resp := &http.Response{
			StatusCode: code,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
		}

		err := validateResponse(resp)
		if err != nil {
			t.Errorf("expected no error for status %d, got %v", code, err)
		}
	}
}

func TestValidateResponse_ErrorWithValidJSON(t *testing.T) {
	body := `{"type":"BadRequest","message":"Invalid input"}`
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/api", nil)

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Request:    req,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}

	err := validateResponse(resp)
	if err == nil || err.Error() != "HTTP POST http://example.com/api failed with status 400: Invalid input" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateResponse_ErrorWithInvalidJSON(t *testing.T) {
	body := `<<<garbage>>>`
	req, _ := http.NewRequest(http.MethodDelete, "http://example.com/delete", nil)

	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Request:    req,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}

	err := validateResponse(resp)
	expectedStart := "HTTP DELETE http://example.com/delete failed with status 500. Raw response:"
	if err == nil || err.Error()[:len(expectedStart)] != expectedStart {
		t.Errorf("unexpected error: %v", err)
	}
}

func createBasicAuthClient(serverURL string) (*RestApiClient, error) {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	return NewRestApiClient(nil, baseURL, "testuser", "testpassword", nil, nil, nil)
}

func createCertAuthClient(serverURL string, serverCACert, clientCert, clientKey []byte) (*RestApiClient, error) {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	return NewRestApiClient(nil, baseURL, "", "", serverCACert, clientCert, clientKey)
}
