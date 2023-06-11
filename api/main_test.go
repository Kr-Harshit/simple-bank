package api

import (
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	DATABASE_URL = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}
