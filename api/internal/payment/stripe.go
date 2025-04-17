package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pl3lee/restjson/internal/database"
	"github.com/redis/go-redis/v9"
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
func (cfg *PaymentConfig) GetSubscriptionStatus(ctx context.Context, customerId string) (subscriptionData, error) {
	cacheKey := fmt.Sprintf("stripe:customer:%s", customerId)
	var subData subscriptionData

	// 1. Check cache
	cachedResult, err := cfg.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - attempt to unmarshal
		if err := json.Unmarshal([]byte(cachedResult), &subData); err == nil {
			return subData, nil // Successfully retrieved from cache
		}
		// Log unmarshal error and proceed to fetch
		fmt.Printf("GetSubscriptionStatus: Failed to unmarshal cached data for customer %s: %v\n", customerId, err)
	} else if err != redis.Nil {
		// Log Redis error (other than cache miss) and proceed to fetch
		fmt.Printf("GetSubscriptionStatus: Redis error for customer %s: %v\n", customerId, err)
	}

	// 2. Cache miss or error - sync from Stripe and update cache
	subData = cfg.syncStripeDataToKV(ctx, customerId) // This function already updates the cache

	return subData, nil
}

func (cfg *PaymentConfig) UpdateSubscriptionStatus(ctx context.Context, customerId string, subscribed bool) (database.User, error) {
	updatedCustomer, err := cfg.Db.UpdateCustomerSubscriptionStatus(ctx, database.UpdateCustomerSubscriptionStatusParams{
		StripeCustomerID: customerId,
		Subscribed:       subscribed,
	})
	if err != nil {
		return database.User{}, fmt.Errorf("UpdateSubscriptionStatus: cannot update subscription status in database: %w", err)
	}
	return updatedCustomer, nil
}
