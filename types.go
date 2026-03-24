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

// ListResponse wraps a paginated list of results.
type ListResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// Meta contains metadata returned with every API response.
type Meta struct {
	RequestID  string  `json:"request_id"`
	Mode       string  `json:"mode"`
	Page       int     `json:"page,omitempty"`
	PerPage    int     `json:"per_page,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
	TotalPages int     `json:"total_pages,omitempty"`
	HasMore    *bool   `json:"has_more,omitempty"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

// apiResponse is the standard envelope for single-object responses.
type apiResponse[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

// apiListResponse is the standard envelope for list responses.
type apiListResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
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
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Status  string `json:"status,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
}

// ListTransactionsParams holds the query parameters for listing transactions.
type ListTransactionsParams struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Type    string `json:"type,omitempty"`
	Status  string `json:"status,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
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
