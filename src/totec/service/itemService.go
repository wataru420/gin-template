package service
import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"totec/models"
	"strings"
)

var itemDao = &models.ItemDao{}

type ItemService struct {}

func (*ItemService) GetEndpoint(c *gin.Context)  {
	id := c.Param("id")

	item, err := itemDao.Get(id)
	if err != nil {
		log.Fatal("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"itemId": item.Id,
		"itemNo": item.No,
		"itemSupplier": item.Supplier,
		"itemSoldQuantity": item.SoldQuantity,
		"itemSalePrice": item.SalePrice,
		"itemTags": item.Tags,
		"itemImage": item.Image,
	})
}

func (*ItemService) GetWebEndpoint(c *gin.Context)  {
	id := c.Param("id")

	item , err := itemDao.Get(id)
	if err != nil {
		log.Println("error")
		c.String(http.StatusInternalServerError, err.Error())
	}

	users, _ := userDao.FindByPostItemId(id,8)
	posts, _ := postDao.FindByPostItemId(id,8)

	c.HTML(http.StatusOK, "itemDetail.tmpl", gin.H{
		"title": "Main website",
		"item": item,
		"users": users,
		"posts": posts,
		"tags": strings.Split(item.Tags,","),
	})
}

//func (*ItemService) FindByPostUserId(id string, limit int) []Item {
//	kvar res = []Item{}

//}