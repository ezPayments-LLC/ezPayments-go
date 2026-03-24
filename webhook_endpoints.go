package ezpayments

import (
	"context"
	"fmt"
)

// WebhookEndpointsResource provides methods for interacting with the webhook endpoints API.
type WebhookEndpointsResource struct {
	client *httpClient
}

// Create creates a new webhook endpoint.
//
//	endpoint, err := client.WebhookEndpoints.Create(ctx, &ezpayments.CreateWebhookEndpointParams{
//	    URL:    "https://example.com/webhooks",
//	    Events: []string{"payment_link.paid", "payment_link.expired"},
//	})
func (r *WebhookEndpointsResource) Create(ctx context.Context, params *CreateWebhookEndpointParams) (*WebhookEndpoint, error) {
	var resp apiResponse[WebhookEndpoint]
	err := r.client.post(ctx, "/webhook-endpoints/", params, &resp, withIdempotencyKey(params.IdempotencyKey))
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a webhook endpoint by ID.
func (r *WebhookEndpointsResource) Get(ctx context.Context, id string) (*WebhookEndpoint, error) {
	var resp apiResponse[WebhookEndpoint]
	err := r.client.get(ctx, fmt.Sprintf("/webhook-endpoints/%s/", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List retrieves all webhook endpoints.
func (r *WebhookEndpointsResource) List(ctx context.Context) (*ListResponse[WebhookEndpoint], error) {
	var resp apiListResponse[WebhookEndpoint]
	err := r.client.get(ctx, "/webhook-endpoints/", &resp)
	if err != nil {
		return nil, err
	}
	return &ListResponse[WebhookEndpoint]{Data: resp.Data, Meta: resp.Meta}, nil
}

// Update updates an existing webhook endpoint.
func (r *WebhookEndpointsResource) Update(ctx context.Context, id string, params *UpdateWebhookEndpointParams) (*WebhookEndpoint, error) {
	var resp apiResponse[WebhookEndpoint]
	err := r.client.patch(ctx, fmt.Sprintf("/webhook-endpoints/%s/", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a webhook endpoint by ID.
func (r *WebhookEndpointsResource) Delete(ctx context.Context, id string) error {
	return r.client.del(ctx, fmt.Sprintf("/webhook-endpoints/%s/", id))
}
