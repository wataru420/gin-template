package service

import (
	"log"
	"net/http"
	"strings"
	"totec/models"

	"github.com/gin-gonic/gin"
)

var userDao = &models.UserDao{}

type UserService struct{}

func (*UserService) GetEndpoint(c *gin.Context) {

	id := c.Param("id")

	user, err := userDao.Get(id)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":         user.Id,
		"userNo":         user.No,
		"userPublicSore": user.PublicScore,
		"userFriends":    strings.Split(user.Friends, ","),
		"userImage":      user.Image,
	})
}

func (*UserService) GetWebEndpoint(c *gin.Context) {
	id := c.Param("id")

	user, err := userDao.Get(id)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	var friends = []models.User{}
	for _, fid := range strings.Split(user.Friends, ",") {
		friend, _ := userDao.Get(fid)
		friends = append(friends, friend)
	}
/*
	items, err := itemDao.FindByPostUserId(id, 8)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}
*/
	posts, err := postDao.FindByPostUserId(id, 8)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}
	var postItems = []models.Item{}
/*
	for _, post := range posts {
		item, _ := itemDao.Get(post.ItemId)
		postItems = append(postItems, item)
	}
*/
	c.HTML(http.StatusOK, "userDetail.tmpl", gin.H{
		"title":        "Main website",
		"user":         user,
		"friendsCount": len(friends),
		"friends":      friends[0:3],
		//"items":        items,
		"posts":        posts,
		"postImages":   postItems,
	})
}

func (*UserService) ListEndpoint(c *gin.Context) {

	type resUser struct {
		UserId string	`json:"userId"`
		UserNo int		`json:"userNo"`
		UserPublicScore int	`json:"userPublicScore"`
		UserFriends []string	`json:"userFriends"`
		UserImage string	`json:"userImage"`
	}
	type res struct {
		Result bool        `json:"result"`
		Data   []resUser `json:"data"`
	}
	limitParam := c.Query("limit")
	limit := "100"
	if limitParam != "" {
		limit = limitParam
	}

	userList, _ := userDao.FindByParam(c, limit)

	var resUserList = []resUser{}
	for _, user := range userList {
		log.Println(user.Id)
		resUserRow := resUser{}
		resUserRow.UserId = user.Id
		resUserRow.UserNo=user.No
		resUserRow.UserPublicScore=user.PublicScore
		resUserRow.UserFriends=strings.Split(user.Friends,",")
		resUserRow.UserImage=user.Image
		resUserList = append(resUserList, resUserRow)
	}
	response := res{true, resUserList}
	log.Println(response)
	c.JSON(http.StatusOK, response)

}
