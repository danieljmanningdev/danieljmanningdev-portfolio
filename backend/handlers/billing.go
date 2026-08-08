package handlers

import (
	"net/http"
	"os"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/services"
)

func HandleCreateCheckout(w http.ResponseWriter, r *http.Request) {
	priceID := os.Getenv("STRIPE_PRICE_ID") // Your recurring product price ID
	domain := "http://localhost:8080"

	checkoutURL, err := services.CreateCheckoutSession(
		priceID,
		domain+"/success?session_id={CHECKOUT_SESSION_ID}",
		domain+"/cancel",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, checkoutURL, http.StatusSeeOther)
}
