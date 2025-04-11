package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
	"github.com/pl3lee/restjson/internal/database"
	"github.com/pl3lee/restjson/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
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
