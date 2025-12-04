// Package payments jest odpowiedzialne za obsługę płatności stripe
//
// Uwaga:
// Tworząc checkoutSession podajemy customerID, ponieważ 1 user może mięc więcej niż 1 produktów (formularzy)
// Flow:
// 1. User chce dodać nowy formularz
// 2. submit -> Post/add z id formularza i emailem
// 3. sprawdzamy czy email istnieje już w bazie i czy ma jakiś formularz:
// 4. jeśli tak to ma już stripe customer id więc musimy z tamtego formularza pobrać to id
// 5. jeśli nie to musimy wygenerować przez funkcje CreateCustomer i wtedy CreatePayment
// (zapisz do bazy danych nowe subscription ID)
// 6. User dokonuje płatności i ląduje na stronie /success
// 7. Webhookami nasłuchujemy stripe'a, żeby kontrolować status is_paid w bazie:
package payments

import (
	"github.com/machayka/mail-service/config"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/customer"
)

type PaymentClient interface {
	CreatePayment(key, priceID, domain string)
}

type Payment struct {
	key     string
	priceID string
	domain  string
}

func NewPaymentClient(cfg *config.Config) *Payment {
	return &Payment{key: cfg.Stripe.Key, priceID: cfg.Stripe.PriceID, domain: cfg.Stripe.Domain}
}

func (p Payment) CreateCustomer(email string) (customerID string, err error) {
	stripe.Key = p.key
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	result, err := customer.New(params)
	return result.ID, err
}

func (p Payment) CreatePayment(customerID, formID string) (checkoutURL string, err error) {
	stripe.Key = p.key
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(p.priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Customer:   stripe.String(customerID),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(p.domain + "/success.html"),
	}
	params.AddMetadata("form_id", formID)
	result, err := session.New(params)
	return result.URL, err
}
