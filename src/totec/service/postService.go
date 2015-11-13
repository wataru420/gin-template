package service
import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"totec/models"
	"strings"
)

var postDao = &models.PostDao{}

type PostService struct {}

func (*PostService) GetEndpoint(c *gin.Context)  {
	id := c.Param("id")

	post, err := postDao.Get(id)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"postId": post.Id,
		"postNo": post.DateTime,
		"postUserId": post.UserId,
		"postItemId": post.ItemId,
		"postItemScore": post.ItemScore,
		"postItemState": post.ItemState,
		"postLikeUsers": post.LikeUsers,
		"postTags": post.Tags,
	})
}

func (*PostService) GetWebEndpoint(c *gin.Context)  {
	id := c.Param("id")

	post , err := postDao.Get(id)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	item , _ := itemDao.Get(post.ItemId)

	var likeUsers = []models.User{}
	for _, fid := range strings.Split(post.LikeUsers,",") {
		friend , _ := userDao.Get(fid)
		likeUsers = append(likeUsers,friend)
	}

	posts, _ := postDao.FindByPostUserId(post.UserId, 8)
	var postItems = []models.Item{}
	for _, post := range posts {
		item, _ := itemDao.Get(post.ItemId)
		postItems = append(postItems, item)
	}


	c.HTML(http.StatusOK, "itemDetail.tmpl", gin.H{
		"title": "Main website",
		"post": post,
		"item": item,
		"likeUsers": likeUsers,
		"postImages": postItems,
		"tags": strings.Split(item.Tags,","),
	})
}
