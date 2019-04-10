package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestQPSControl(t *testing.T) {
	engine := gin.Default()
	engine.Use(QPSControl(map[string]int64{"/api": 1000}))
	engine.GET("/api", func(context *gin.Context) {
		fmt.Println(context.Request.RemoteAddr)
	})
	engine.Run(":8080")
}
