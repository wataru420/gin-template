package controllers
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"totec/service"
)

var userService = &service.UserService{}

func InitRooter(e *gin.Engine) {
	e.GET("/status", func(c *gin.Context) {c.String(http.StatusOK, "ok")})

	json := e.Group("/json")
	{
		json.GET("user/detail/:id", auth(), userService.GetEndpoint)
		json.GET("user/list", auth(), userService.ListEndpoint)
	}

	web := e.Group("/web")
	{
		web.GET("user/detail/:id", auth(), userService.GetWebEndpoint)
		web.GET("user/list", auth(), userService.ListWebEndpoint)
	}

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/join/
	e.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		service.ApiAuth(c)
	}
}
