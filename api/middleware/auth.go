package middleware

import (
	"net/http"
	"strings"

	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/ProstoyVadila/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func abortUnauthorized(c *gin.Context, err e.ErrUnauthorized) {
	c.AbortWithStatusJSON(err.StatusCode(), errorResponse(err))
}

// Auth is a middleware that checks the Authorization header
func Auth(tokenMaker token.Maker) gin.HandlerFunc {
	log.Info().Msg("Setting Auth middleware")
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if authorizationHeader == "" {
			err := e.ErrUnauthorized{Msg: "authorization token is not provided"}
			abortUnauthorized(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := e.ErrUnauthorized{Msg: "invalid authorization header"}
			abortUnauthorized(ctx, err)
			return
		}
		authorizationType := fields[0]
		if authorizationType != AuthorizationTypeBearer {
			err := e.ErrUnauthorized{Msg: "unsupported authorization type:", Obj: authorizationType}
			abortUnauthorized(ctx, err)
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.ValidateToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
