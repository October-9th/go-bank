package api

import (
	"os"
	"testing"
	"time"

	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestSever(t *testing.T, store sqlc.Store) *Server {
	config := util.Config{
		TokenSymmectricKey:  util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server

}
func TestMain(m *testing.M) {
	gin.SetMode(gin.DebugMode)

	os.Exit(m.Run())
}
