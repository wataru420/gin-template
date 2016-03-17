package service

import (
	"log"
	"net/http"
	"strings"
	"time"
	"totec/models"

	"github.com/gin-gonic/gin"
)

var postDao = &models.PostDao{}

type PostService struct{}

func (*PostService) GetEndpoint(c *gin.Context) {
	id := c.Param("id")

	post, err := postDao.Get(id)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"postId":        post.Id,
		"postNo":        post.DateTime,
		"postUserId":    post.UserId,
		"postItemId":    post.ItemId,
		"postItemScore": post.ItemScore,
		"postItemState": post.ItemState,
		"postLikeUsers": post.LikeUsers,
		"postTags":      post.Tags,
	})
}

func (*PostService) GetWebEndpoint(c *gin.Context) {
	id := c.Param("id")

	post, err := postDao.Get(id)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	item, _ := itemDao.Get(post.ItemId)

	var likeUsers = []models.User{}
	for _, fid := range strings.Split(post.LikeUsers, ",") {
		friend, _ := userDao.Get(fid)
		likeUsers = append(likeUsers, friend)
	}

	posts, _ := postDao.FindByPostUserId(post.UserId, 8)
	var postItems = []models.Item{}
	for _, p := range posts {
		item, _ := itemDao.Get(p.ItemId)
		postItems = append(postItems, item)
	}

	var userItems = []models.User{}
	userPosts, _ := postDao.FindByPostItemId(post.ItemId, 8)
	for _, p := range userPosts {
		item, _ := userDao.Get(p.UserId)
		userItems = append(userItems, item)
	}

	user, _ := userDao.Get(post.UserId)

	c.HTML(http.StatusOK, "postDetail.tmpl", gin.H{
		"title":         "Main website",
		"post":          post,
		"item":          item,
		"user":          user,
		"postTime":      time.Unix(int64(post.DateTime), 0).Format("2006年1月2日 15:4"),
		"likeUserCount": len(likeUsers),
		"likeUsers":     likeUsers[0:3],
		"postImages":    postItems,
		"userImages":    userItems,
		//"tags":          strings.Split(item.Tags, ","),
	})
}


func (*PostService) ListEndpoint(c *gin.Context) {

	type resItem struct {
		Id string	`json:"postId"`
		DateTime int		`json:"postDateTime"`
		UserId string	`json:"postUserId"`
		ItemId string	`json:"postItemId"`
		ItemScore int	`json:"postItemScore"`
		ItemState string	`json:"postItemState"`
		LikeUsers []string	`json:"postLikeUsers"`
		Tags []string	`json:"postTags"`
	}
	type res struct {
		Result bool        `json:"result"`
		Data   []resItem `json:"data"`
	}
	limitParam := c.Query("limit")
	limit := "100"
	if limitParam != "" {
		limit = limitParam
	}

	itemList, _ := postDao.FindByParam(c, limit)

	var resItemList = []resItem{}
	for _, row := range itemList {
		resItemRow := resItem{}
		resItemRow.Id = row.Id
		resItemRow.DateTime = row.DateTime
		resItemRow.UserId = row.UserId
		resItemRow.ItemId = row.ItemId
		resItemRow.ItemScore = row.ItemScore
		resItemRow.ItemState = row.ItemState
		resItemRow.LikeUsers = strings.Split(row.LikeUsers,",")
		resItemRow.Tags = strings.Split(row.Tags,",")
		resItemList = append(resItemList, resItemRow)
	}
	response := res{true, resItemList}
	log.Println(response)
	c.JSON(http.StatusOK, response)

}
