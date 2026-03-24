# ezpayments-go

Go SDK for the [ezPayments](https://app.ezpayments.co) Merchant API v3.

[![Go Reference](https://pkg.go.dev/badge/github.com/elkhayyat/ezpayments-go.svg)](https://pkg.go.dev/github.com/elkhayyat/ezpayments-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/elkhayyat/ezpayments-go)](https://goreportcard.com/report/github.com/elkhayyat/ezpayments-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features

- Zero external dependencies (stdlib only)
- Full coverage of the ezPayments Merchant API v3
- Webhook signature verification (HMAC-SHA256)
- Functional options pattern for client configuration
- Idempotent request support
- Proper Go error handling with typed errors
- Context support on all methods

## Requirements

- Go 1.21 or later

## Installation

```bash
go get github.com/elkhayyat/ezpayments-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    ezpayments "github.com/elkhayyat/ezpayments-go"
)

func main() {
    client := ezpayments.New("sk_live_your_api_key")

    link, err := client.PaymentLinks.Create(context.Background(), &ezpayments.CreatePaymentLinkParams{
        Amount:      "50.00",
        Description: "Invoice #1234",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment link: %s\n", link.URL)
}
```

## Client Configuration

```go
// Default configuration
client := ezpayments.New("sk_live_xxx")

// Custom base URL (for sandbox)
client := ezpayments.New("sk_test_xxx",
    ezpayments.WithBaseURL("https://sandbox.ezpayments.co"),
)

// Custom timeout
client := ezpayments.New("sk_live_xxx",
    ezpayments.WithTimeout(60 * time.Second),
)

// Custom HTTP client
client := ezpayments.New("sk_live_xxx",
    ezpayments.WithHTTPClient(&http.Client{
        Transport: customTransport,
    }),
)
```

## API Reference

### Payment Links

```go
// Create a payment link
link, err := client.PaymentLinks.Create(ctx, &ezpayments.CreatePaymentLinkParams{
    Amount:          "50.00",
    Description:     "Invoice #1234",
    CustomerName:    "John Doe",
    CustomerEmail:   "john@example.com",
    ReferenceNumber: "INV-1234",
    IdempotencyKey:  "unique-request-id",  // optional
})

// Get a payment link
link, err := client.PaymentLinks.Get(ctx, "pl_abc123")

// List payment links
links, err := client.PaymentLinks.List(ctx, &ezpayments.ListPaymentLinksParams{
    Page:   1,
    PerPage: 20,
    Status: "active",
})

// Update a payment link
amount := "75.00"
link, err := client.PaymentLinks.Update(ctx, "pl_abc123", &ezpayments.UpdatePaymentLinkParams{
    Amount: &amount,
})

// Delete a payment link
err := client.PaymentLinks.Delete(ctx, "pl_abc123")

// Get fee breakdown
fees, err := client.PaymentLinks.GetFees(ctx, "pl_abc123")
fmt.Printf("Fee: %s, Total: %s\n", fees.Fee, fees.Total)
```

### Transactions

```go
// Get a transaction
txn, err := client.Transactions.Get(ctx, "txn_abc123")

// List transactions
txns, err := client.Transactions.List(ctx, &ezpayments.ListTransactionsParams{
    Page:    1,
    PerPage: 20,
    Status:  "completed",
    Type:    "payment",
})
```

### Webhook Endpoints

```go
// Create a webhook endpoint
endpoint, err := client.WebhookEndpoints.Create(ctx, &ezpayments.CreateWebhookEndpointParams{
    URL:    "https://example.com/webhooks",
    Events: []string{"payment_link.paid", "payment_link.expired"},
})
// Save endpoint.Secret for verifying webhook signatures

// Get a webhook endpoint
endpoint, err := client.WebhookEndpoints.Get(ctx, "we_abc123")

// List webhook endpoints
endpoints, err := client.WebhookEndpoints.List(ctx)

// Update a webhook endpoint
newURL := "https://example.com/new-webhook"
endpoint, err := client.WebhookEndpoints.Update(ctx, "we_abc123", &ezpayments.UpdateWebhookEndpointParams{
    URL: &newURL,
})

// Delete a webhook endpoint
err := client.WebhookEndpoints.Delete(ctx, "we_abc123")
```

### API Keys

```go
// Create an API key (full key only returned on creation)
key, err := client.APIKeys.Create(ctx, &ezpayments.CreateAPIKeyParams{
    Name: "Production Key",
})
fmt.Printf("Key: %s\n", *key.Key)  // Save this, shown only once

// List API keys
keys, err := client.APIKeys.List(ctx)

// Delete (revoke) an API key
err := client.APIKeys.Delete(ctx, "key_abc123")
```

## Webhook Verification

Verify incoming webhook signatures to ensure they originated from ezPayments:

```go
import (
    "io"
    "net/http"

    ezpayments "github.com/elkhayyat/ezpayments-go"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    sig := r.Header.Get("X-EzPayments-Signature")
    if err := ezpayments.VerifyWebhookSignature("whsec_your_secret", sig, body); err != nil {
        http.Error(w, "invalid signature", http.StatusForbidden)
        return
    }

    // Signature is valid, process the event
    w.WriteHeader(http.StatusOK)
}
```

The signature header format is `t=timestamp,v1=hmac_hex`, where the signature is computed as `HMAC-SHA256(secret, "{timestamp}.{raw_body}")`.

## Error Handling

All API errors are returned as `*ezpayments.APIError`, which can be inspected using `errors.As`:

```go
import "errors"

link, err := client.PaymentLinks.Get(ctx, "pl_nonexistent")
if err != nil {
    var apiErr *ezpayments.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Type:    %s\n", apiErr.Type)       // e.g. "invalid_request_error"
        fmt.Printf("Message: %s\n", apiErr.Message)    // e.g. "Payment link not found"
        fmt.Printf("Code:    %s\n", apiErr.Code)       // e.g. "resource_missing"
        fmt.Printf("Status:  %d\n", apiErr.StatusCode) // e.g. 404

        switch apiErr.StatusCode {
        case 401:
            // Invalid API key
        case 404:
            // Resource not found
        case 429:
            // Rate limited
        }
    }
    // Non-API errors (network issues, timeouts, etc.)
    log.Fatal(err)
}
```

## Response Metadata

All list responses include metadata:

```go
links, err := client.PaymentLinks.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Request ID: %s\n", links.Meta.RequestID)
fmt.Printf("Mode: %s\n", links.Meta.Mode)  // "live" or "test"
fmt.Printf("Total: %d\n", links.Meta.TotalCount)
fmt.Printf("Page: %d\n", links.Meta.Page)
```

## License

MIT License - see [LICENSE](LICENSE) for details.
