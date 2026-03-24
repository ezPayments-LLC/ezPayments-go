package ezpayments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidSignatureHeader is returned when the signature header format is invalid.
	ErrInvalidSignatureHeader = errors.New("ezpayments: invalid signature header format")

	// ErrSignatureMismatch is returned when the computed signature does not match.
	ErrSignatureMismatch = errors.New("ezpayments: signature verification failed")

	// ErrMissingTimestamp is returned when the signature header has no timestamp.
	ErrMissingTimestamp = errors.New("ezpayments: missing timestamp in signature header")

	// ErrMissingSignature is returned when the signature header has no v1 signature.
	ErrMissingSignature = errors.New("ezpayments: missing v1 signature in header")
)

// VerifyWebhookSignature verifies the signature of an incoming webhook request.
// The signature header has the format: t=timestamp,v1=hmac_hex
// The expected signature is HMAC-SHA256(secret, "{timestamp}.{rawBody}").
//
//	body, _ := io.ReadAll(r.Body)
//	sig := r.Header.Get("X-EzPayments-Signature")
//	if err := ezpayments.VerifyWebhookSignature("whsec_xxx", sig, body); err != nil {
//	    http.Error(w, "invalid signature", http.StatusForbidden)
//	    return
//	}
func VerifyWebhookSignature(secret string, signatureHeader string, body []byte) error {
	if signatureHeader == "" {
		return ErrInvalidSignatureHeader
	}

	var timestamp, signature string

	parts := strings.Split(signatureHeader, ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signature = kv[1]
		}
	}

	if timestamp == "" {
		return ErrMissingTimestamp
	}
	if signature == "" {
		return ErrMissingSignature
	}

	payload := fmt.Sprintf("%s.%s", timestamp, string(body))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrSignatureMismatch
	}

	return nil
}
