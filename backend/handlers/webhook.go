package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

func StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	// If a webhook secret is set, verify the signature
	if endpointSecret != "" {
		signatureHeader := r.Header.Get("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
		if err != nil {
			http.Error(w, "Webhook signature verification failed", http.StatusBadRequest)
			return
		}

		if event.Type == "checkout.session.completed" {
			var session stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &session)
			if err == nil {
				// TODO: Update your database here to mark the client's invoice/project as paid!
				// e.g., session.CustomerEmail or session.ID
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
