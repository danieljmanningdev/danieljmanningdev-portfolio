package handlers

import (
	"fmt"
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

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Base build price
	basePrice := 2500
	addonTotal := 0

	// Check selected features from checkboxes
	features := r.Form["feature"]
	for _, f := range features {
		switch f {
		case "messaging":
			addonTotal += 1000
		case "billing":
			addonTotal += 1500
		}
	}

	total := basePrice + addonTotal

	// Return just the dynamic results HTML fragment for HTMX to swap
	html := fmt.Sprintf(`
		<div id="calc-result" class="border border-zinc-800 bg-zinc-950 p-6 rounded-lg flex flex-col justify-between">
			<div>
				<h3 class="text-sm font-mono text-zinc-400 uppercase tracking-wider mb-4">Estimated Summary</h3>
				<div class="space-y-3 mb-6">
					<div class="flex justify-between text-sm">
						<span class="text-zinc-400">Base Build:</span>
						<span class="font-mono text-white">£%d</span>
					</div>
					<div class="flex justify-between text-sm">
						<span class="text-zinc-400">Selected Add-ons:</span>
						<span class="font-mono text-white">£%d</span>
					</div>
					<div class="border-t border-zinc-800 pt-3 flex justify-between font-bold text-lg">
						<span class="text-white">Total (Excl. VAT):</span>
						<span class="font-mono text-cyan-400">£%d</span>
					</div>
				</div>
			</div>
			<button class="w-full bg-cyan-600 text-white font-medium py-2.5 rounded hover:bg-cyan-500 transition">
				Lock In This Scope (£%d)
			</button>
		</div>
	`, basePrice, addonTotal, total, total)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Stub for Stripe checkout redirect once they lock it in
func StripeCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	// We will hook up stripe-go session creation here next
	w.Write([]byte("Redirecting to Stripe Checkout..."))
}
