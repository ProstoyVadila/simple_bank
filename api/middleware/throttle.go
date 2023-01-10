package middleware

import (
	"net/http"

	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func Throttling(maxEventsPerSec int, maxBurstSize int) gin.HandlerFunc {
	log.Info().Msg("Setting Throttling middleware")
	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(context *gin.Context) {
		if limiter.Allow() {
			context.Next()
			return
		}

		context.Error(e.ErrThrottling{})
		context.AbortWithStatus(http.StatusTooManyRequests)
	}
}
