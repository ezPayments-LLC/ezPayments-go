// Package main demonstrates how to use the ezPayments Go SDK.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	ezpayments "github.com/elkhayyat/ezpayments-go"
)

func main() {
	// Create a client with your API key
	client := ezpayments.New("sk_live_your_api_key",
		ezpayments.WithBaseURL("https://app.ezpayments.co"),
	)

	ctx := context.Background()

	// --- Create a payment link ---
	link, err := client.PaymentLinks.Create(ctx, &ezpayments.CreatePaymentLinkParams{
		Amount:        "50.00",
		Description:   "Invoice #1234",
		CustomerName:  "John Doe",
		CustomerEmail: "john@example.com",
	})
	if err != nil {
		var apiErr *ezpayments.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("API error: %s (code=%s, status=%d)", apiErr.Message, apiErr.Code, apiErr.StatusCode)
		}
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Created payment link: %s\n", link.URL)

	// --- List payment links ---
	links, err := client.PaymentLinks.List(ctx, &ezpayments.ListPaymentLinksParams{
		Page:   1,
		Status: "active",
	})
	if err != nil {
		log.Fatalf("Error listing links: %v", err)
	}
	fmt.Printf("Found %d payment links\n", len(links.Data))

	// --- Get fees for a payment link ---
	fees, err := client.PaymentLinks.GetFees(ctx, link.ID)
	if err != nil {
		log.Fatalf("Error getting fees: %v", err)
	}
	fmt.Printf("Fee: %s, Total: %s\n", fees.Fee, fees.Total)

	// --- List transactions ---
	txns, err := client.Transactions.List(ctx, &ezpayments.ListTransactionsParams{
		PerPage: 10,
	})
	if err != nil {
		log.Fatalf("Error listing transactions: %v", err)
	}
	for _, txn := range txns.Data {
		fmt.Printf("Transaction %s: %s %s (%s)\n", txn.ID, txn.Amount, txn.Currency, txn.Status)
	}

	// --- Create a webhook endpoint ---
	endpoint, err := client.WebhookEndpoints.Create(ctx, &ezpayments.CreateWebhookEndpointParams{
		URL:    "https://example.com/webhooks/ezpayments",
		Events: []string{"payment_link.paid", "payment_link.expired"},
	})
	if err != nil {
		log.Fatalf("Error creating webhook: %v", err)
	}
	fmt.Printf("Webhook endpoint created: %s (secret: %s)\n", endpoint.ID, endpoint.Secret)

	// --- Webhook handler example ---
	webhookSecret := endpoint.Secret

	http.HandleFunc("/webhooks/ezpayments", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		sig := r.Header.Get("X-EzPayments-Signature")
		if err := ezpayments.VerifyWebhookSignature(webhookSecret, sig, body); err != nil {
			http.Error(w, "invalid signature", http.StatusForbidden)
			return
		}

		// Process the webhook event
		fmt.Println("Verified webhook received")
		w.WriteHeader(http.StatusOK)
	})
}
