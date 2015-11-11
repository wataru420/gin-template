package service
import (
	"log"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"totec/models"
)

var userDao = &models.UserDao{}

type UserService struct {}

func (*UserService) GetEndpoint(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, id + " is not int")
		return
	}
	user, err := userDao.Get(i)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.Id,
		"name": user.Name,
		"type": user.Type,
	})
}

func (*UserService) GetWebEndpoint(c *gin.Context)  {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, id + " is not int")
		return
	}
	user , err := userDao.Get(i)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.HTML(http.StatusOK, "userDetail.tmpl", gin.H{
		"title": "Main website",
		"user": user,
	})
}
func (*UserService) ListEndpoint(c *gin.Context) {
	type res struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Type int `json:"type"`
	}

	resList := []res{}

	userList, err := userDao.GetList()
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	for _, user := range userList {
		r := res{Id:user.Id,Name:user.Name,Type:user.Type}
		resList = append(resList, r)
	}

	c.JSON(http.StatusOK, resList)
}
