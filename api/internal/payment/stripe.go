package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	var result []*stripe.Subscription = []*stripe.Subscription{}
	result = subscription.List(&stripe.SubscriptionListParams{
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
		SubscriptionId: subscription.ID,
		Status:         string(subscription.Status),
	}

	// Safely access subscription item data
	if subscription.Items != nil && len(subscription.Items.Data) > 0 {
		item := subscription.Items.Data[0]
		if item != nil {
			subData.CurrentPeriodStart = int(item.CurrentPeriodStart)
			subData.CurrentPeriodEnd = int(item.CurrentPeriodEnd)
			if item.Price != nil {
				subData.PriceId = item.Price.ID
			}
		}
	}

	// Safely access default payment method data
	if subscription.DefaultPaymentMethod != nil && subscription.DefaultPaymentMethod.Card != nil {
		subData.PaymentMethod = paymentMethod{
			Brand: string(subscription.DefaultPaymentMethod.Card.Brand),
			Last4: subscription.DefaultPaymentMethod.Card.Last4,
		}
	}
	fmt.Print(subData)
	jsonData, err := json.Marshal(subData)
	if err != nil {
		log.Printf("syncStripeDataToKV: error marshalling subData to JSON: %v\n", err)
	}
	err = cfg.Rdb.Set(ctx, fmt.Sprintf("stripe:customer:%s", customerId), jsonData, 0).Err()
	if err != nil {
		log.Printf("syncStripeDataToKV: error in saving jsonData to redis: %v\n", err)
	}
	return subData
}

// TODO: GetSubscriptionStatusFromKV
// retrives subscription status for customer in redis
func (cfg *PaymentConfig) GetSubscriptionStatusFromKV(ctx context.Context, customerId string) (subscriptionData, bool, error) {
	cacheKey := fmt.Sprintf("stripe:customer:%s", customerId)
	var subData subscriptionData

	cachedResult, err := cfg.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - attempt to unmarshal
		if err := json.Unmarshal([]byte(cachedResult), &subData); err == nil {
			return subData, true, nil // Successfully retrieved from cache
		}
		return subData, false, fmt.Errorf("GetSubscriptionStatusFromKV: error in unmarshalling data %w", err)
	} else if err != redis.Nil {
		return subData, false, fmt.Errorf("GetSubscriptionStatusFromKV: redis error for customer: %v", err)
	} else {
		return subData, false, fmt.Errorf("GetSubscriptionStatusFromKV: not in KV: %v", err)
	}
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
