package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/service-kit/qps-controler/qps"
	"net/http"
	"sync"
)

var (
	qpsControler *qps.QPSControler
	once         sync.Once
)

func init() {
	qpsControler = new(qps.QPSControler)
	qpsControler.Init()
}

func QPSControl(apiQPSLimitMap map[string]int64) gin.HandlerFunc {
	for key, val := range apiQPSLimitMap {
		qpsControler.NewRule(key, qps.Rule{Limit: val})
	}
	return func(c *gin.Context) {
		if !qpsControler.Pass(c.Request.URL.Path, nil) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  "fail",
				"error": "request later",
			})
			c.Abort()
			return
		}
	}
}
