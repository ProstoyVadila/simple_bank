package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Register(router *gin.Engine) {
	router.Use(
		RecoveryMiddleware(),
		DefaultLoggerMiddleware(),
		CORSMiddleware(),
	)
}

func RecoveryMiddleware() gin.HandlerFunc {
	log.Info().Msg("Setting Recovery middleware")
	return gin.Recovery()
}
