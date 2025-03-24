package ratelimit

import (
	"net/http"
	"time"

	"github.com/pl3lee/restjson/internal/utils"
	"github.com/redis/go-redis/v9"
)

// TokenBucketRateLimiter is a middleware that implements a token bucket rate limiting algorithm.
// It limits the number of requests a client can make within a certain time frame.
//
// Parameters:
// - rdb: Redis client used to store and retrieve rate limiting data.
// - capacity: Maximum number of tokens (requests) a client can have at any given time.
// - refillRate: Rate at which tokens are added to the bucket (tokens per second).
// - expiration: Time in seconds after which the rate limiting data expires in Redis.
//
// The middleware works as follows:
// 1. It identifies the client by their IP address.
// 2. It retrieves the current number of tokens and the last access time from Redis.
// 3. It refills the tokens based on the elapsed time since the last access.
// 4. If the client has enough tokens, it allows the request and consumes one token.
// 5. If the client does not have enough tokens, it responds with a 429 Too Many Requests status.
func TokenBucketRateLimiter(rdb *redis.Client, capacity int, refillRate float64, expiration int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			clientIP := r.RemoteAddr
			tokensKey := "rate_limit:" + clientIP + ":tokens"
			lastAccessKey := "rate_limit:" + clientIP + ":last_access"

			// get current tokens
			currentTokens, err := rdb.Get(ctx, tokensKey).Int()
			if err == redis.Nil {
				// create if does not exist
				currentTokens = capacity
				rdb.Set(ctx, tokensKey, capacity, time.Duration(expiration)*time.Second)
			} else if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "rate limiting error", err)
				return
			}

			// get last access time
			lastAccessStr, err := rdb.Get(ctx, lastAccessKey).Result()
			var lastAccess time.Time
			if err == redis.Nil {
				// if not found, then default to 1 hour ago
				lastAccess = time.Now().Add(-1 * time.Hour)
				rdb.Set(ctx, lastAccessKey, time.Now().Format(time.RFC3339), time.Duration(expiration)*time.Second)
			} else if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "rate limiting error", err)
				return
			} else {
				lastAccess, _ = time.Parse(time.RFC3339, lastAccessStr)
			}

			// token refill
			now := time.Now()
			elapsed := now.Sub(lastAccess).Seconds()
			tokensToAdd := int(elapsed * refillRate)
			newTokens := currentTokens + tokensToAdd
			if newTokens > capacity {
				newTokens = capacity
			}

			// check if request can be processed
			if newTokens < 1 {
				utils.RespondWithError(w, http.StatusTooManyRequests, "rate limit exceeded", nil)
				return
			}

			// consume token
			newTokens--

			// update redis
			rdb.Set(ctx, tokensKey, newTokens, time.Duration(expiration)*time.Second)
			rdb.Set(ctx, lastAccessKey, now.Format(time.RFC3339), time.Duration(expiration)*time.Second)

			// next request
			next.ServeHTTP(w, r)
		})
	}
}
