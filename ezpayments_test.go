package ezpayments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestServer creates a test HTTP server and an ezpayments client pointing to it.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := New("sk_test_123", WithBaseURL(server.URL))
	t.Cleanup(server.Close)
	return client, server
}

// jsonResponse writes a JSON response with the given status code.
func jsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func TestNew(t *testing.T) {
	client := New("sk_live_test")
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.PaymentLinks == nil {
		t.Fatal("expected PaymentLinks resource")
	}
	if client.Transactions == nil {
		t.Fatal("expected Transactions resource")
	}
	if client.WebhookEndpoints == nil {
		t.Fatal("expected WebhookEndpoints resource")
	}
	if client.APIKeys == nil {
		t.Fatal("expected APIKeys resource")
	}
}

func TestNewWithOptions(t *testing.T) {
	client := New("sk_live_test",
		WithBaseURL("https://custom.example.com"),
	)
	if client.httpClient.baseURL != "https://custom.example.com" {
		t.Fatalf("expected custom base URL, got %s", client.httpClient.baseURL)
	}
}

func TestAuthorizationHeader(t *testing.T) {
	var gotAuth string
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		jsonResponse(w, http.StatusOK, apiResponse[PaymentLink]{
			Data: PaymentLink{ID: "pl_1"},
			Meta: Meta{RequestID: "req_1", Mode: "test"},
		})
	})

	_, err := client.PaymentLinks.Get(context.Background(), "pl_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Bearer sk_test_123"
	if gotAuth != expected {
		t.Fatalf("expected auth header %q, got %q", expected, gotAuth)
	}
}

func TestUserAgentHeader(t *testing.T) {
	var gotUA string
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		jsonResponse(w, http.StatusOK, apiResponse[PaymentLink]{
			Data: PaymentLink{ID: "pl_1"},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	_, err := client.PaymentLinks.Get(context.Background(), "pl_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "ezpayments-go/" + Version
	if gotUA != expected {
		t.Fatalf("expected user-agent %q, got %q", expected, gotUA)
	}
}

// --- Payment Links ---

func TestPaymentLinksCreate(t *testing.T) {
	var gotBody map[string]interface{}
	var gotIdemKey string
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v3/payment-links/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotIdemKey = r.Header.Get("Idempotency-Key")
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		jsonResponse(w, http.StatusCreated, apiResponse[PaymentLink]{
			Data: PaymentLink{
				ID:          "pl_abc123",
				Amount:      "50.00",
				Description: "Test payment",
				Status:      "active",
			},
			Meta: Meta{RequestID: "req_1", Mode: "test"},
		})
	})

	link, err := client.PaymentLinks.Create(context.Background(), &CreatePaymentLinkParams{
		Amount:         "50.00",
		Description:    "Test payment",
		IdempotencyKey: "idem_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if link.ID != "pl_abc123" {
		t.Fatalf("expected ID pl_abc123, got %s", link.ID)
	}
	if link.Amount != "50.00" {
		t.Fatalf("expected amount 50.00, got %s", link.Amount)
	}
	if gotIdemKey != "idem_123" {
		t.Fatalf("expected idempotency key idem_123, got %s", gotIdemKey)
	}
	if gotBody["amount"] != "50.00" {
		t.Fatalf("expected request body amount 50.00, got %v", gotBody["amount"])
	}
}

func TestPaymentLinksGet(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v3/payment-links/pl_123/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		jsonResponse(w, http.StatusOK, apiResponse[PaymentLink]{
			Data: PaymentLink{ID: "pl_123", Amount: "25.00"},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	link, err := client.PaymentLinks.Get(context.Background(), "pl_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.ID != "pl_123" {
		t.Fatalf("expected ID pl_123, got %s", link.ID)
	}
}

func TestPaymentLinksList(t *testing.T) {
	nextURL := "https://app.ezpayments.co/api/v3/payment-links/?cursor=pl_2"
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("status") != "active" {
			t.Fatalf("expected status=active, got %s", r.URL.Query().Get("status"))
		}
		jsonResponse(w, http.StatusOK, apiListResponse[PaymentLink]{
			Data: apiListData[PaymentLink]{
				Results: []PaymentLink{
					{ID: "pl_1", Amount: "10.00"},
					{ID: "pl_2", Amount: "20.00"},
				},
				Next:     &nextURL,
				Previous: nil,
			},
			Meta: Meta{RequestID: "req_1", Mode: "test"},
		})
	})

	resp, err := client.PaymentLinks.List(context.Background(), &ListPaymentLinksParams{
		ListParams: ListParams{Limit: 10},
		Status:     "active",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 links, got %d", len(resp.Results))
	}
	if !resp.HasMore() {
		t.Fatal("expected HasMore() to be true")
	}
	if resp.Previous != nil {
		t.Fatal("expected Previous to be nil")
	}
}

func TestPaymentLinksListWithStartingAfter(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("starting_after") != "pl_abc" {
			t.Fatalf("expected starting_after=pl_abc, got %s", r.URL.Query().Get("starting_after"))
		}
		jsonResponse(w, http.StatusOK, apiListResponse[PaymentLink]{
			Data: apiListData[PaymentLink]{
				Results: []PaymentLink{{ID: "pl_def"}},
			},
			Meta: Meta{RequestID: "req_2"},
		})
	})

	resp, err := client.PaymentLinks.List(context.Background(), &ListPaymentLinksParams{
		ListParams: ListParams{StartingAfter: "pl_abc"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 link, got %d", len(resp.Results))
	}
	if resp.HasMore() {
		t.Fatal("expected HasMore() to be false when Next is nil")
	}
}

func TestPaymentLinksListNilParams(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %s", r.URL.RawQuery)
		}
		jsonResponse(w, http.StatusOK, apiListResponse[PaymentLink]{
			Data: apiListData[PaymentLink]{
				Results: []PaymentLink{{ID: "pl_1"}},
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	resp, err := client.PaymentLinks.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 link, got %d", len(resp.Results))
	}
}

func TestPaymentLinksUpdate(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		jsonResponse(w, http.StatusOK, apiResponse[PaymentLink]{
			Data: PaymentLink{ID: "pl_1", Amount: "75.00"},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	newAmount := "75.00"
	link, err := client.PaymentLinks.Update(context.Background(), "pl_1", &UpdatePaymentLinkParams{
		Amount: &newAmount,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link.Amount != "75.00" {
		t.Fatalf("expected amount 75.00, got %s", link.Amount)
	}
}

func TestPaymentLinksDelete(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.PaymentLinks.Delete(context.Background(), "pl_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPaymentLinksGetFees(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/payment-links/pl_1/fees/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		jsonResponse(w, http.StatusOK, apiResponse[PaymentLinkFees]{
			Data: PaymentLinkFees{
				Amount:   "50.00",
				Fee:      "1.50",
				Total:    "51.50",
				Currency: "USD",
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	fees, err := client.PaymentLinks.GetFees(context.Background(), "pl_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fees.Fee != "1.50" {
		t.Fatalf("expected fee 1.50, got %s", fees.Fee)
	}
}

// --- Transactions ---

func TestTransactionsGet(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/transactions/txn_1/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		jsonResponse(w, http.StatusOK, apiResponse[Transaction]{
			Data: Transaction{ID: "txn_1", Amount: "50.00", Status: "completed"},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	txn, err := client.Transactions.Get(context.Background(), "txn_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if txn.ID != "txn_1" {
		t.Fatalf("expected ID txn_1, got %s", txn.ID)
	}
}

func TestTransactionsList(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusOK, apiListResponse[Transaction]{
			Data: apiListData[Transaction]{
				Results: []Transaction{{ID: "txn_1"}, {ID: "txn_2"}},
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	resp, err := client.Transactions.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(resp.Results))
	}
}

func TestTransactionsListWithPagination(t *testing.T) {
	nextURL := "https://app.ezpayments.co/api/v3/transactions/?cursor=txn_2"
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "5" {
			t.Fatalf("expected limit=5, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("status") != "completed" {
			t.Fatalf("expected status=completed, got %s", r.URL.Query().Get("status"))
		}
		jsonResponse(w, http.StatusOK, apiListResponse[Transaction]{
			Data: apiListData[Transaction]{
				Results: []Transaction{{ID: "txn_1"}},
				Next:    &nextURL,
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	resp, err := client.Transactions.List(context.Background(), &ListTransactionsParams{
		ListParams: ListParams{Limit: 5},
		Status:     "completed",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(resp.Results))
	}
	if !resp.HasMore() {
		t.Fatal("expected HasMore() to be true")
	}
}

// --- Webhook Endpoints ---

func TestWebhookEndpointsCreate(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		jsonResponse(w, http.StatusCreated, apiResponse[WebhookEndpoint]{
			Data: WebhookEndpoint{
				ID:     "we_1",
				URL:    "https://example.com/webhooks",
				Events: []string{"payment_link.paid"},
				Secret: "whsec_xxx",
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	endpoint, err := client.WebhookEndpoints.Create(context.Background(), &CreateWebhookEndpointParams{
		URL:    "https://example.com/webhooks",
		Events: []string{"payment_link.paid"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if endpoint.ID != "we_1" {
		t.Fatalf("expected ID we_1, got %s", endpoint.ID)
	}
	if endpoint.Secret != "whsec_xxx" {
		t.Fatalf("expected secret whsec_xxx, got %s", endpoint.Secret)
	}
}

func TestWebhookEndpointsList(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusOK, apiListResponse[WebhookEndpoint]{
			Data: apiListData[WebhookEndpoint]{
				Results: []WebhookEndpoint{{ID: "we_1"}, {ID: "we_2"}},
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	resp, err := client.WebhookEndpoints.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(resp.Results))
	}
}

func TestWebhookEndpointsDelete(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.WebhookEndpoints.Delete(context.Background(), "we_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- API Keys ---

func TestAPIKeysCreate(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		key := "sk_live_full_key_value"
		jsonResponse(w, http.StatusCreated, apiResponse[APIKey]{
			Data: APIKey{
				ID:     "key_1",
				Name:   "Production",
				Prefix: "sk_live_",
				Key:    &key,
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	key, err := client.APIKeys.Create(context.Background(), &CreateAPIKeyParams{
		Name: "Production",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.ID != "key_1" {
		t.Fatalf("expected ID key_1, got %s", key.ID)
	}
	if key.Key == nil || *key.Key != "sk_live_full_key_value" {
		t.Fatal("expected key value in response")
	}
}

func TestAPIKeysList(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusOK, apiListResponse[APIKey]{
			Data: apiListData[APIKey]{
				Results: []APIKey{{ID: "key_1", Name: "Prod"}, {ID: "key_2", Name: "Dev"}},
			},
			Meta: Meta{RequestID: "req_1"},
		})
	})

	resp, err := client.APIKeys.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(resp.Results))
	}
}

// --- Pagination ---

func TestListResponseHasMore(t *testing.T) {
	next := "https://app.ezpayments.co/api/v3/payment-links/?cursor=abc"
	withNext := &ListResponse[PaymentLink]{
		Results: []PaymentLink{{ID: "pl_1"}},
		Next:    &next,
	}
	if !withNext.HasMore() {
		t.Fatal("expected HasMore() to be true when Next is set")
	}

	withoutNext := &ListResponse[PaymentLink]{
		Results: []PaymentLink{{ID: "pl_1"}},
		Next:    nil,
	}
	if withoutNext.HasMore() {
		t.Fatal("expected HasMore() to be false when Next is nil")
	}
}

// --- Error Handling ---

func TestAPIError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, http.StatusNotFound, apiErrorEnvelope{
			Error: APIError{
				Type:    "invalid_request_error",
				Message: "Payment link not found",
				Code:    "resource_missing",
			},
		})
	})

	_, err := client.PaymentLinks.Get(context.Background(), "pl_nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Type != "invalid_request_error" {
		t.Fatalf("expected type invalid_request_error, got %s", apiErr.Type)
	}
	if apiErr.Code != "resource_missing" {
		t.Fatalf("expected code resource_missing, got %s", apiErr.Code)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Payment link not found" {
		t.Fatalf("expected message 'Payment link not found', got %s", apiErr.Message)
	}
}

func TestAPIErrorFormat(t *testing.T) {
	err := &APIError{
		Type:       "invalid_request_error",
		Message:    "Amount is required",
		Code:       "missing_param",
		Param:      "amount",
		StatusCode: 400,
	}
	expected := "ezpayments: Amount is required (type=invalid_request_error, code=missing_param, status=400)"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestMalformedErrorResponse(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	_, err := client.PaymentLinks.Get(context.Background(), "pl_1")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Fatalf("expected status 500, got %d", apiErr.StatusCode)
	}
	if apiErr.Type != "unknown_error" {
		t.Fatalf("expected type unknown_error, got %s", apiErr.Type)
	}
}

// --- Webhook Verification ---

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "whsec_test_secret"
	timestamp := "1700000000"
	body := []byte(`{"event":"payment_link.paid","data":{"id":"pl_1"}}`)

	payload := fmt.Sprintf("%s.%s", timestamp, string(body))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))

	header := fmt.Sprintf("t=%s,v1=%s", timestamp, sig)

	err := VerifyWebhookSignature(secret, header, body)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestVerifyWebhookSignatureInvalid(t *testing.T) {
	err := VerifyWebhookSignature("secret", "t=123,v1=invalidsig", []byte("body"))
	if err != ErrSignatureMismatch {
		t.Fatalf("expected ErrSignatureMismatch, got: %v", err)
	}
}

func TestVerifyWebhookSignatureEmptyHeader(t *testing.T) {
	err := VerifyWebhookSignature("secret", "", []byte("body"))
	if err != ErrInvalidSignatureHeader {
		t.Fatalf("expected ErrInvalidSignatureHeader, got: %v", err)
	}
}

func TestVerifyWebhookSignatureMissingTimestamp(t *testing.T) {
	err := VerifyWebhookSignature("secret", "v1=abc", []byte("body"))
	if err != ErrMissingTimestamp {
		t.Fatalf("expected ErrMissingTimestamp, got: %v", err)
	}
}

func TestVerifyWebhookSignatureMissingV1(t *testing.T) {
	err := VerifyWebhookSignature("secret", "t=123", []byte("body"))
	if err != ErrMissingSignature {
		t.Fatalf("expected ErrMissingSignature, got: %v", err)
	}
}

// --- Context Cancellation ---

func TestContextCancellation(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response -- context should cancel before we respond
		<-r.Context().Done()
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.PaymentLinks.Get(ctx, "pl_1")
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}
