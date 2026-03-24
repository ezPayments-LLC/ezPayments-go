package ezpayments

// PaymentLink represents a payment link resource.
type PaymentLink struct {
	ID              string  `json:"id"`
	Amount          string  `json:"amount"`
	Description     string  `json:"description"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email"`
	ReferenceNumber string  `json:"reference_number"`
	Token           string  `json:"token"`
	URL             string  `json:"url"`
	Status          string  `json:"status"`
	ExpiresAt       *string `json:"expires_at"`
	PaidAt          *string `json:"paid_at"`
	ViewCount       int     `json:"view_count"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// PaymentLinkFees represents fee information for a payment link.
type PaymentLinkFees struct {
	Amount        string `json:"amount"`
	Fee           string `json:"fee"`
	Total         string `json:"total"`
	Currency      string `json:"currency"`
	FeePercentage string `json:"fee_percentage"`
}

// Transaction represents a transaction resource.
type Transaction struct {
	ID              string  `json:"id"`
	Amount          string  `json:"amount"`
	Fee             string  `json:"fee"`
	Net             string  `json:"net"`
	Currency        string  `json:"currency"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	Description     string  `json:"description"`
	ReferenceNumber string  `json:"reference_number"`
	PaymentLinkID   *string `json:"payment_link_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// WebhookEndpoint represents a webhook endpoint resource.
type WebhookEndpoint struct {
	ID            string   `json:"id"`
	URL           string   `json:"url"`
	Description   string   `json:"description"`
	Secret        string   `json:"secret"`
	Events        []string `json:"events"`
	Status        string   `json:"status"`
	LastDelivered *string  `json:"last_delivered"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

// APIKey represents an API key resource.
type APIKey struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Prefix    string  `json:"prefix"`
	Key       *string `json:"key"`
	LastUsed  *string `json:"last_used"`
	ExpiresAt *string `json:"expires_at"`
	CreatedAt string  `json:"created_at"`
}

// ListParams holds cursor-based pagination parameters accepted by all list
// endpoints. Both fields are optional.
type ListParams struct {
	// Limit is the maximum number of items to return per page (1-100, default 20).
	Limit int `json:"limit,omitempty"`

	// StartingAfter is a cursor for forward pagination. Pass the ID of the
	// last item from the previous page to fetch the next page of results.
	StartingAfter string `json:"starting_after,omitempty"`
}

// ListResponse wraps a paginated list of results returned by the API.
type ListResponse[T any] struct {
	// Results contains the items for the current page.
	Results []T `json:"results"`

	// Next is the full URL for the next page, or nil if there are no more pages.
	Next *string `json:"next"`

	// Previous is the full URL for the previous page, or nil if this is the first page.
	Previous *string `json:"previous"`

	// Meta contains request metadata (request ID, mode).
	Meta Meta `json:"meta"`
}

// HasMore reports whether there is a next page of results.
func (r *ListResponse[T]) HasMore() bool {
	return r.Next != nil
}

// Meta contains metadata returned with every API response.
type Meta struct {
	RequestID string `json:"request_id"`
	Mode      string `json:"mode"`
}

// apiResponse is the standard envelope for single-object responses.
type apiResponse[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

// apiListData is the nested data payload for list responses containing
// cursor-based pagination fields alongside the result set.
type apiListData[T any] struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

// apiListResponse is the standard envelope for list responses.
type apiListResponse[T any] struct {
	Data apiListData[T] `json:"data"`
	Meta Meta           `json:"meta"`
}

// CreatePaymentLinkParams holds the parameters for creating a payment link.
type CreatePaymentLinkParams struct {
	Amount          string  `json:"amount"`
	Description     string  `json:"description"`
	CustomerName    string  `json:"customer_name,omitempty"`
	CustomerEmail   string  `json:"customer_email,omitempty"`
	ReferenceNumber string  `json:"reference_number,omitempty"`
	ExpiresAt       *string `json:"expires_at,omitempty"`
	IdempotencyKey  string  `json:"-"`
}

// UpdatePaymentLinkParams holds the parameters for updating a payment link.
type UpdatePaymentLinkParams struct {
	Amount          *string `json:"amount,omitempty"`
	Description     *string `json:"description,omitempty"`
	CustomerName    *string `json:"customer_name,omitempty"`
	CustomerEmail   *string `json:"customer_email,omitempty"`
	ReferenceNumber *string `json:"reference_number,omitempty"`
	ExpiresAt       *string `json:"expires_at,omitempty"`
}

// ListPaymentLinksParams holds the query parameters for listing payment links.
type ListPaymentLinksParams struct {
	ListParams
	Status string `json:"status,omitempty"`
}

// ListTransactionsParams holds the query parameters for listing transactions.
type ListTransactionsParams struct {
	ListParams
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}

// ListWebhookEndpointsParams holds the query parameters for listing webhook endpoints.
type ListWebhookEndpointsParams struct {
	ListParams
}

// ListAPIKeysParams holds the query parameters for listing API keys.
type ListAPIKeysParams struct {
	ListParams
}

// CreateWebhookEndpointParams holds the parameters for creating a webhook endpoint.
type CreateWebhookEndpointParams struct {
	URL            string   `json:"url"`
	Description    string   `json:"description,omitempty"`
	Events         []string `json:"events"`
	IdempotencyKey string   `json:"-"`
}

// UpdateWebhookEndpointParams holds the parameters for updating a webhook endpoint.
type UpdateWebhookEndpointParams struct {
	URL         *string  `json:"url,omitempty"`
	Description *string  `json:"description,omitempty"`
	Events      []string `json:"events,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

// CreateAPIKeyParams holds the parameters for creating an API key.
type CreateAPIKeyParams struct {
	Name           string  `json:"name"`
	ExpiresAt      *string `json:"expires_at,omitempty"`
	IdempotencyKey string  `json:"-"`
}
