package payment

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/utils"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/webhook"
)

type checkoutParams struct {
	PriceId string `json:"priceId"`
}

func (cfg *PaymentConfig) HandlerCheckout(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	var checkoutRequest checkoutParams
	if err := utils.DecodeRequest(r, &checkoutRequest); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "cannot decode request", err)
		return
	}

	user, err := cfg.Db.GetUserById(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error getting user from database", err)
		return
	}
	stripeCustomerId := user.StripeCustomerID

	// customer does not exist, create it in stripe
	if stripeCustomerId == "" {
		newCustomer, err := customer.New(&stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Metadata: map[string]string{
				"userId": userId.String(),
			},
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating stripe customer", err)
			return
		}
		// store in db
		stripeCustomerId = newCustomer.ID
		_, err = cfg.Db.UpdateCustomerId(r.Context(), database.UpdateCustomerIdParams{
			ID:               user.ID,
			StripeCustomerID: stripeCustomerId,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error updating stripe customer", err)
			return
		}
	}

	// create checkout session
	checkoutParams := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerId),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(checkoutRequest.PriceId),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(cfg.ClientURL + "/success"),
		CancelURL:  stripe.String(cfg.ClientURL + "/cancel"),
	}

	s, err := session.New(checkoutParams)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create stripe checkout session", err)
		return
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}

// called from frontend /success page
func (cfg *PaymentConfig) HandlerSuccess(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIDContextKey).(uuid.UUID)
	user, err := cfg.Db.GetUserById(r.Context(), userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "cannot get user", err)
		return
	}
	stripeCustomerId := user.StripeCustomerID
	if stripeCustomerId == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "cannot find customer id", err)
		return
	}

	subscriptionData := cfg.syncStripeDataToKV(r.Context(), stripeCustomerId)
	utils.RespondWithJSON(w, http.StatusOK, subscriptionData)
}

var allowedEvents map[stripe.EventType]bool = map[stripe.EventType]bool{
	stripe.EventTypeCheckoutSessionCompleted:                 true,
	stripe.EventTypeCustomerSubscriptionCreated:              true,
	stripe.EventTypeCustomerSubscriptionUpdated:              true,
	stripe.EventTypeCustomerSubscriptionDeleted:              true,
	stripe.EventTypeCustomerSubscriptionPaused:               true,
	stripe.EventTypeCustomerSubscriptionResumed:              true,
	stripe.EventTypeCustomerSubscriptionPendingUpdateApplied: true,
	stripe.EventTypeCustomerSubscriptionPendingUpdateExpired: true,
	stripe.EventTypeCustomerSubscriptionTrialWillEnd:         true,
	stripe.EventTypeInvoicePaid:                              true,
	stripe.EventTypeInvoicePaymentFailed:                     true,
	stripe.EventTypeInvoicePaymentActionRequired:             true,
	stripe.EventTypeInvoiceUpcoming:                          true,
	stripe.EventTypeInvoiceMarkedUncollectible:               true,
	stripe.EventTypeInvoicePaymentSucceeded:                  true,
	stripe.EventTypePaymentIntentSucceeded:                   true,
	stripe.EventTypePaymentIntentPaymentFailed:               true,
	stripe.EventTypePaymentIntentCanceled:                    true,
}

func (cfg *PaymentConfig) HandlerStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error reading request body", err)
		return
	}

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := cfg.StripeWebhookSecret
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "error verifying webhook signature", err)
		return
	}
	if !allowedEvents[event.Type] {
		utils.RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	customerId := ""
	if customerObj, ok := event.Data.Object["customer"]; ok {
		if id, ok := customerObj.(string); ok {
			customerId = id
		}
	}

	if customerId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "customer ID not found in event", nil)
		return
	}

	subscriptionData := cfg.syncStripeDataToKV(r.Context(), customerId)
	fmt.Printf("got subscription data: %v", subscriptionData)
	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
