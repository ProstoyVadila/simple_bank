package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/ProstoyVadila/simple_bank/token"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type testServer struct {
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

func newTestServer(t *testing.T) *testServer {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	router := gin.New()
	tokenMaker, err := token.NewPaseto(config.TokenSymmetricKey)
	require.NoError(t, err)
	return &testServer{
		router:     router,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}
