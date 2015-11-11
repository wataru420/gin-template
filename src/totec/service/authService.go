package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"octo/models"
)

	var appDao = &models.AppDAO{}

func ApiAuth(c * gin.Context) {
	var app = models.App{}
	err := appDao.Get(&app, 1)
	if err != nil {
		fmt.Println(err)
	}
	c.Set("app",app)
	val, _ := c.Get("app")
	if str, ok := val.(models.App); ok {
		println(str.AppName)
	}
}