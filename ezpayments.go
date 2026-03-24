// Package ezpayments provides a Go client for the ezPayments Merchant API v3.
//
// Create a client with your API key and start making requests:
//
//	client := ezpayments.New("sk_live_xxx")
//	link, err := client.PaymentLinks.Create(ctx, &ezpayments.CreatePaymentLinkParams{
//	    Amount:      "50.00",
//	    Description: "Invoice #1234",
//	})
package ezpayments

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default API base URL.
	DefaultBaseURL = "https://app.ezpayments.co"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// Version is the SDK version.
	Version = "0.1.0"

	apiVersion = "v3"
)

// Client is the ezPayments API client.
type Client struct {
	// PaymentLinks provides access to the payment links API.
	PaymentLinks *PaymentLinksResource

	// Transactions provides access to the transactions API.
	Transactions *TransactionsResource

	// WebhookEndpoints provides access to the webhook endpoints API.
	WebhookEndpoints *WebhookEndpointsResource

	// APIKeys provides access to the API keys management API.
	APIKeys *APIKeysResource

	httpClient *httpClient
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.httpClient.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom net/http client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient.client = hc
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.client.Timeout = timeout
	}
}

// New creates a new ezPayments API client with the given API key.
// Use functional options to customize the client:
//
//	client := ezpayments.New("sk_live_xxx",
//	    ezpayments.WithBaseURL("https://sandbox.ezpayments.co"),
//	    ezpayments.WithTimeout(60 * time.Second),
//	)
func New(apiKey string, opts ...ClientOption) *Client {
	hc := &httpClient{
		baseURL: DefaultBaseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	c := &Client{
		httpClient: hc,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.PaymentLinks = &PaymentLinksResource{client: hc}
	c.Transactions = &TransactionsResource{client: hc}
	c.WebhookEndpoints = &WebhookEndpointsResource{client: hc}
	c.APIKeys = &APIKeysResource{client: hc}

	return c
}
