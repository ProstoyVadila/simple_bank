package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}
