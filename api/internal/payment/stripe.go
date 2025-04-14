package payment

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/subscription"
)

type paymentMethod struct {
	Brand string `json:"brand"`
	Last4 string `json:"last4"`
}

type subscriptionData struct {
	SubscriptionId     string        `json:"subscriptionId"`
	Status             string        `json:"status"`
	PriceId            string        `json:"priceId"`
	CurrentPeriodStart int           `json:"currentPeriodStart"`
	CurrentPeriodEnd   int           `json:"currentPeriodEnd"`
	PaymentMethod      paymentMethod `json:"paymentMethod"`
}

func (cfg *PaymentConfig) syncStripeDataToKV(ctx context.Context, customerId string) subscriptionData {
	result := subscription.List(&stripe.SubscriptionListParams{
		Customer: stripe.String(customerId),
		Status:   stripe.String("all"),
		ListParams: stripe.ListParams{
			Limit: stripe.Int64(1),
		},
	}).SubscriptionList().Data
	if len(result) == 0 {
		subData := subscriptionData{
			Status: "none",
		}
		cfg.Rdb.Set(ctx, fmt.Sprintf("stripe:customer:%s", customerId), subData, 0)
		return subData
	}

	subscription := result[0]

	subData := subscriptionData{
		SubscriptionId:     subscription.ID,
		Status:             string(subscription.Status),
		PriceId:            subscription.Items.Data[0].Price.ID,
		CurrentPeriodStart: int(subscription.Items.Data[0].CurrentPeriodStart),
		CurrentPeriodEnd:   int(subscription.Items.Data[0].CurrentPeriodEnd),
		PaymentMethod: paymentMethod{
			Brand: string(subscription.DefaultPaymentMethod.Card.Brand),
			Last4: subscription.DefaultPaymentMethod.Card.Last4,
		},
	}
	cfg.Rdb.Set(ctx, fmt.Sprintf("stripe:customer:%s", customerId), subData, 0)
	return subData
}

// TODO: GetSubscriptionStatus
// retrives subscription status for customer
// first checks KV, if not found, trigger sync from stripe and check redis again
// if still not found, then no active subscription
