package ezpayments

import (
	"context"
	"fmt"
	"net/url"
)

// APIKeysResource provides methods for interacting with the API keys management API.
type APIKeysResource struct {
	client *httpClient
}

// Create creates a new API key. The full key value is only available in the
// response to this call and cannot be retrieved again.
func (r *APIKeysResource) Create(ctx context.Context, params *CreateAPIKeyParams) (*APIKey, error) {
	var resp apiResponse[APIKey]
	err := r.client.post(ctx, "/api-keys/", params, &resp, withIdempotencyKey(params.IdempotencyKey))
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List retrieves a paginated list of API keys. Key values are not included in
// list responses.
func (r *APIKeysResource) List(ctx context.Context, params *ListAPIKeysParams) (*ListResponse[APIKey], error) {
	qp := url.Values{}
	if params != nil {
		qp = encodeListParams(params.ListParams, nil)
	}
	var resp apiListResponse[APIKey]
	err := r.client.getWithQuery(ctx, "/api-keys/", qp, &resp)
	if err != nil {
		return nil, err
	}
	return &ListResponse[APIKey]{
		Results:  resp.Data.Results,
		Next:     resp.Data.Next,
		Previous: resp.Data.Previous,
		Meta:     resp.Meta,
	}, nil
}

// Delete revokes an API key by ID.
func (r *APIKeysResource) Delete(ctx context.Context, id string) error {
	return r.client.del(ctx, fmt.Sprintf("/api-keys/%s/", id))
}
