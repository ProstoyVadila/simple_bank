package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ProstoyVadila/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

func abortUnauthorized(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
}

// Auth is a middleware that checks the Authorization header
func Auth(tokenMaker token.Maker) gin.HandlerFunc {
	log.Info().Msg("Setting Auth middleware")
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if authorizationHeader == "" {
			err := errors.New("authorization token is not provided")
			abortUnauthorized(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authorization header")
			abortUnauthorized(ctx, err)
			return
		}
		authorizationType := fields[0]
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %v", authorizationType)
			abortUnauthorized(ctx, err)
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.ValidateToken(accessToken)
		if err != nil {
			abortUnauthorized(ctx, err)
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
