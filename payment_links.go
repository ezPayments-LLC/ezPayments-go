package ezpayments

import (
	"context"
	"fmt"
	"net/url"
)

// PaymentLinksResource provides methods for interacting with the payment links API.
type PaymentLinksResource struct {
	client *httpClient
}

// Create creates a new payment link.
//
//	link, err := client.PaymentLinks.Create(ctx, &ezpayments.CreatePaymentLinkParams{
//	    Amount:      "50.00",
//	    Description: "Invoice #1234",
//	})
func (r *PaymentLinksResource) Create(ctx context.Context, params *CreatePaymentLinkParams) (*PaymentLink, error) {
	var resp apiResponse[PaymentLink]
	err := r.client.post(ctx, "/payment-links/", params, &resp, withIdempotencyKey(params.IdempotencyKey))
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a payment link by ID.
func (r *PaymentLinksResource) Get(ctx context.Context, id string) (*PaymentLink, error) {
	var resp apiResponse[PaymentLink]
	err := r.client.get(ctx, fmt.Sprintf("/payment-links/%s/", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List retrieves a paginated list of payment links.
func (r *PaymentLinksResource) List(ctx context.Context, params *ListPaymentLinksParams) (*ListResponse[PaymentLink], error) {
	qp := url.Values{}
	if params != nil {
		qp = encodeListParams(params.ListParams, map[string]string{
			"status": params.Status,
		})
	}
	var resp apiListResponse[PaymentLink]
	err := r.client.getWithQuery(ctx, "/payment-links/", qp, &resp)
	if err != nil {
		return nil, err
	}
	return &ListResponse[PaymentLink]{
		Results:  resp.Data.Results,
		Next:     resp.Data.Next,
		Previous: resp.Data.Previous,
		Meta:     resp.Meta,
	}, nil
}

// Update updates an existing payment link.
func (r *PaymentLinksResource) Update(ctx context.Context, id string, params *UpdatePaymentLinkParams) (*PaymentLink, error) {
	var resp apiResponse[PaymentLink]
	err := r.client.patch(ctx, fmt.Sprintf("/payment-links/%s/", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a payment link by ID.
func (r *PaymentLinksResource) Delete(ctx context.Context, id string) error {
	return r.client.del(ctx, fmt.Sprintf("/payment-links/%s/", id))
}

// GetFees retrieves the fee breakdown for a payment link.
func (r *PaymentLinksResource) GetFees(ctx context.Context, id string) (*PaymentLinkFees, error) {
	var resp apiResponse[PaymentLinkFees]
	err := r.client.get(ctx, fmt.Sprintf("/payment-links/%s/fees/", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
