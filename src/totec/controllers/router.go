package controllers
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"octo/service"
)

func InitRooter(e *gin.Engine) {
	e.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	v1 := e.Group("/v1")
	{
		v1.GET("upload", auth(), service.Upload)
		v1.GET("list/:version/:revision", auth(), service.ListEndpoint)
	}

}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		service.ApiAuth(c)
	}
}
