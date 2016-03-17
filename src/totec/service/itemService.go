package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"totec/models"
)

type ItemService struct{}

var itemDao = &models.ItemDao{}

func (*ItemService) ReadItemEndpoint(c *gin.Context) {

	type data struct {
		Id    string   `json:"itemId"`
		Name  string   `json:"itemName"`
		Type  string   `json:"itemType"`
		Value  int   `json:"itemValue"`
		EffectTarget string   `json:"itemEffectTarget"`
		EffectValue  int   `json:"itemEffectValue"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	id := c.Query("targetItemId")

	item, _ := itemDao.Get(id)

	var list = []data{}
	row := data{}
	row.Id = item.Id
	row.Name = item.Name
	row.Type = item.Type
	row.Value = item.Value
	row.EffectTarget = item.EffectTarget
	row.EffectValue = item.EffectValue
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func (*ItemService) UpdateItemEndpoint(c *gin.Context) {

	type data struct {
		Id    string   `json:"itemId"`
		Name  string   `json:"itemName"`
		Type  string   `json:"itemType"`
		Value  int   `json:"itemValue"`
		EffectTarget string   `json:"itemEffectTarget"`
		EffectValue  int   `json:"itemEffectValue"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	id := c.Query("targetItemId")

	value := c.Query("newItemValue")

	itemDao.Update(id, value)

	item, _ := itemDao.Get(id)

	var list = []data{}
	row := data{}
	row.Id = item.Id
	row.Name = item.Name
	row.Type = item.Type
	row.Value = item.Value
	row.EffectTarget = item.EffectTarget
	row.EffectValue = item.EffectValue
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}