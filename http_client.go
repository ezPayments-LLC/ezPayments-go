package ezpayments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// httpClient wraps net/http with authentication and JSON handling.
type httpClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// request executes an HTTP request and decodes the JSON response into result.
func (h *httpClient) request(ctx context.Context, method, path string, body interface{}, result interface{}, opts ...requestOption) error {
	reqURL := h.baseURL + "/api/" + apiVersion + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("ezpayments: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("ezpayments: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ezpayments-go/"+Version)

	for _, opt := range opts {
		opt(req)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("ezpayments: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ezpayments: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(respBody, resp.StatusCode)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("ezpayments: failed to decode response: %w", err)
		}
	}

	return nil
}

// get performs a GET request.
func (h *httpClient) get(ctx context.Context, path string, result interface{}) error {
	return h.request(ctx, http.MethodGet, path, nil, result)
}

// getWithQuery performs a GET request with query parameters.
func (h *httpClient) getWithQuery(ctx context.Context, path string, params url.Values, result interface{}) error {
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}
	return h.request(ctx, http.MethodGet, path, nil, result)
}

// post performs a POST request.
func (h *httpClient) post(ctx context.Context, path string, body interface{}, result interface{}, opts ...requestOption) error {
	return h.request(ctx, http.MethodPost, path, body, result, opts...)
}

// patch performs a PATCH request.
func (h *httpClient) patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.request(ctx, http.MethodPatch, path, body, result)
}

// del performs a DELETE request.
func (h *httpClient) del(ctx context.Context, path string) error {
	return h.request(ctx, http.MethodDelete, path, nil, nil)
}

// requestOption modifies an HTTP request before it is sent.
type requestOption func(*http.Request)

// withIdempotencyKey sets the Idempotency-Key header.
func withIdempotencyKey(key string) requestOption {
	return func(req *http.Request) {
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
	}
}

// parseAPIError attempts to parse an API error from a response body.
func parseAPIError(body []byte, statusCode int) error {
	var envelope apiErrorEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return &APIError{
			Type:       "unknown_error",
			Message:    fmt.Sprintf("unexpected status %d: %s", statusCode, string(body)),
			StatusCode: statusCode,
		}
	}
	envelope.Error.StatusCode = statusCode
	return &envelope.Error
}

// encodeListParams converts common list parameters to url.Values.
func encodeListParams(page, perPage int, extra map[string]string) url.Values {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Set("per_page", strconv.Itoa(perPage))
	}
	for k, v := range extra {
		if v != "" {
			params.Set(k, v)
		}
	}
	return params
}
