package api

import (
	"github.com/gin-gonic/gin"
)

func registerMiddlewares(router *gin.Engine) {
	router.Use(ErrorMiddleware())

}

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
