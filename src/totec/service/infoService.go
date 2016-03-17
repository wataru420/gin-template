package service
import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)


type InfoService struct {}

func (*InfoService) GetInfoEndpoint(c *gin.Context)  {

	type resItem struct {
		Information string	`json:"information"`
	}
	type res struct {
		Result bool        `json:"result"`
		Data   []resItem `json:"data"`
	}


	var resItemList = []resItem{}
	resItemRow := resItem{}
	resItemRow.Information = "本日のお知らせ"
	resItemList = append(resItemList, resItemRow)

	response := res{true, resItemList}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}
