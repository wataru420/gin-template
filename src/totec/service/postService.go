package service
import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"totec/models"
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

	c.HTML(http.StatusOK, "itemDetail.tmpl", gin.H{
		"title": "Main website",
		"post": post,
	})
}
