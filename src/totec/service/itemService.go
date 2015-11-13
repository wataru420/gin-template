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
	var postItems = []models.User{}
	for _, post := range posts {
		item, _ := userDao.Get(post.UserId)
		postItems = append(postItems, item)
	}

	c.HTML(http.StatusOK, "itemDetail.tmpl", gin.H{
		"title": "Main website",
		"item": item,
		"users": users,
		"posts": posts,
		"tags": strings.Split(item.Tags,","),
		"postImages": postItems,
	})
}


func (*ItemService) ListEndpoint(c *gin.Context) {

	type resItem struct {
		Id string	`json:"itemId"`
		No int		`json:"itemNo"`
		Supplier string	`json:"itemSoldQuantity"`
		SoldQuantity int	`json:"itemSoldQuantity"`
		SalePrice int	`json:"itemSalePrice"`
		Tags []string	`json:"itemTags"`
		Image string	`json:"itemImage"`
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

	itemList, _ := itemDao.FindByParam(c, limit)

	var resItemList = []resItem{}
	for _, row := range itemList {
		resItemRow := resItem{}
		resItemRow.Id = row.Id
		resItemRow.No = row.No
		resItemRow.Supplier = row.Supplier
		resItemRow.SoldQuantity = row.SoldQuantity
		resItemRow.SalePrice = row.SalePrice
		resItemRow.Tags = strings.Split(row.Tags,",")
		resItemRow.Image = row.Image
		resItemList = append(resItemList, resItemRow)
	}
	response := res{true, resItemList}
	log.Println(response)
	c.JSON(http.StatusOK, response)

}
//func (*ItemService) FindByPostUserId(id string, limit int) []Item {
//	kvar res = []Item{}

//}