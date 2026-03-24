package ezpayments

import (
	"context"
	"fmt"
	"net/url"
)

// TransactionsResource provides methods for interacting with the transactions API.
type TransactionsResource struct {
	client *httpClient
}

// Get retrieves a transaction by ID.
func (r *TransactionsResource) Get(ctx context.Context, id string) (*Transaction, error) {
	var resp apiResponse[Transaction]
	err := r.client.get(ctx, fmt.Sprintf("/transactions/%s/", id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List retrieves a paginated list of transactions.
func (r *TransactionsResource) List(ctx context.Context, params *ListTransactionsParams) (*ListResponse[Transaction], error) {
	qp := url.Values{}
	if params != nil {
		qp = encodeListParams(params.ListParams, map[string]string{
			"type":   params.Type,
			"status": params.Status,
		})
	}
	var resp apiListResponse[Transaction]
	err := r.client.getWithQuery(ctx, "/transactions/", qp, &resp)
	if err != nil {
		return nil, err
	}
	return &ListResponse[Transaction]{
		Results:  resp.Data.Results,
		Next:     resp.Data.Next,
		Previous: resp.Data.Previous,
		Meta:     resp.Meta,
	}, nil
}
