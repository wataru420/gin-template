package service
import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"totec/models"
	"strings"
	"strconv"
)

var userDao = &models.UserDao{}

type UserService struct {}

func (*UserService) GetEndpoint(c *gin.Context)  {


	id := c.Param("id")

	user, err := userDao.Get(id)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": user.Id,
		"userNo": user.No,
		"userPublicSore": user.PublicScore,
		"userFriends": strings.Split(user.Friends,","),
		"userImage": user.Image,
	})
}

func (*UserService) GetWebEndpoint(c *gin.Context)  {
	id := c.Param("id")

	user , err := userDao.Get(id)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	var friends = []models.User{}
	for _, fid := range strings.Split(user.Friends,",") {
		friend , _ := userDao.Get(fid)
		friends = append(friends,friend)
	}

	items, err := itemDao.FindByPostUserId(id,8)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	posts, err := postDao.FindByPostUserId(id,8)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	var postItems = []models.Item{}
	for _, post := range posts {
		item, _ := itemDao.Get(post.ItemId)
		postItems = append(postItems, item)
	}

	c.HTML(http.StatusOK, "userDetail.tmpl", gin.H{
		"title": "Main website",
		"user": user,
		"friends": friends,
		"items": items,
		"posts": posts,
		"postImages": postItems,
	})
}

func (*UserService) ListEndpoint(c *gin.Context) {
	type res struct {
		Result string `json:"result"`
		Data []models.User `json:"name"`
	}
	limitParam := c.Param("limit")
	limit := 100
	if limitParam == "" {
		limit,_ = strconv.Atoi(limitParam)
	}
	log.Println(limit)

	param := c.Param("findByPostId")
	if param != "" {
		var userList = []models.User{}
		user, _ := userDao.Get(param)
		userList = append(userList, user)
		c.JSON(http.StatusOK,res{"true",userList})

	}
//	resList := []res{}
//
//	userList, err := userDao.GetList()
//	if err != nil {
//		log.Fatal("error")
//		c.String(http.StatusInternalServerError, err.Error())
//	}
//
//	for _, user := range userList {
//		r := res{Id:user.Id,Name:user.Name,Type:user.Type}
//		resList = append(resList, r)
//	}
//
//	c.JSON(http.StatusOK, resList)
//}
//
//func (*UserService) ListWebEndpoint(c *gin.Context) {
//	type res struct {
//		Id int `json:"id"`
//		Name string `json:"name"`
//		Type int `json:"type"`
//	}
//
//	resList := []res{}
//
//	userList, err := userDao.GetList()
//	if err != nil {
//		log.Fatal("error")
//		c.String(http.StatusInternalServerError, err.Error())
//	}
//
//	for _, user := range userList {
//		r := res{Id:user.Id,Name:user.Name,Type:user.Type}
//		resList = append(resList, r)
//	}
//
//	c.HTML(http.StatusOK, "userList.tmpl", gin.H{
//		"title": "Main website",
//		"userList": resList,
//	})
}
