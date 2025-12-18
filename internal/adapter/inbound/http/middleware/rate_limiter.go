package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// DDoS / Bot Protection
func BotProtectionMiddleware(store fiber.Storage, maxReqs int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               maxReqs,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		Storage:           store,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == ""
		},
	})
}

// Normal Request (done by users)
func RateLimitMiddleware(store fiber.Storage, maxReqs int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               maxReqs,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		Storage:           store,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID, ok := c.Locals("user_id").(uint)
			if !ok {
				return ""
			}
			return strconv.Itoa(int(userID))
		},
		Next: func(c *fiber.Ctx) bool {
			return c.Locals("user_id") == nil
		},
	})
}
