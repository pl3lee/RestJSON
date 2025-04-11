package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pl3lee/restjson/internal/auth"
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

	stripeCustomerId, err := cfg.Rdb.Get(r.Context(), fmt.Sprintf("stripe:user:%s", userId)).Result()
	if err == nil {
		// cache hit, customer exists
	} else if err == redis.Nil {
		// cache miss, no customer exists
		// create stripe customer
		newCustomer, err := customer.New(&stripe.CustomerParams{
			Email: stripe.String(userId.String()),
			Metadata: map[string]string{
				"userId": userId.String(),
			},
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error creating stripe customer", err)
			return
		}
		// store in redis
		stripeCustomerId = newCustomer.ID
		err = cfg.Rdb.Set(r.Context(), fmt.Sprintf("stripe:user:%s", userId), stripeCustomerId, 0).Err()
		if err != nil {
			log.Printf("failed to cache Stripe customer ID for user %s: %v", userId, err)
		}

	} else {
		// actual error
		utils.RespondWithError(w, http.StatusInternalServerError, "error checking redis for stripe customer", err)
		return
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
