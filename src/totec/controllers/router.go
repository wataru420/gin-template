package controllers
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"totec/service"
)

var userService = &service.UserService{}
var itemService = &service.ItemService{}
var postService = &service.PostService{}

func InitRooter(e *gin.Engine) {
	e.GET("/status", func(c *gin.Context) {c.String(http.StatusOK, "ok")})

	e.GET("/searchUser", userService.GetEndpoint)
	e.GET("/searchItem", func(c *gin.Context) {c.String(http.StatusOK, "ok")})
	e.GET("/searchPost", func(c *gin.Context) {c.String(http.StatusOK, "ok")})

	e.GET("/user/:id", userService.GetWebEndpoint)
	e.GET("/item/:id", itemService.GetWebEndpoint)
	e.GET("/post/:id", postService.GetWebEndpoint)

	e.Static("/static", "./static")
//
//	json := e.Group("/json")
//	{
//		json.GET("user/detail/:id", auth(), userService.GetEndpoint)
//		json.GET("user/list", auth(), userService.ListEndpoint)
//	}
//
//	web := e.Group("/web")
//	{
//		web.GET("user/detail/:id", auth(), userService.GetWebEndpoint)
//		web.GET("user/list", auth(), userService.ListWebEndpoint)
//	}
//
//	// However, this one will match /user/john/ and also /user/john/send
//	// If no other routers match /user/john, it will redirect to /user/join/
//	e.GET("/user/:name/*action", func(c *gin.Context) {
//		name := c.Param("name")
//		action := c.Param("action")
//		message := name + " is " + action
//		c.String(http.StatusOK, message)
//	})

}

