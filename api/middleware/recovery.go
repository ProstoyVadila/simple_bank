package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Recovery() gin.HandlerFunc {
	log.Info().Msg("Setting Recovery middleware")
	return gin.Recovery()
}
