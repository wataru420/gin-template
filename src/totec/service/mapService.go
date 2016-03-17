package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"totec/models"
)

type MapService struct{}

var mapDao = &models.MapDao{}

func (*MapService) ReadMapEndpoint(c *gin.Context) {

	type data struct {
		Id    string   `json:"mapId"`
		Name  string   `json:"mapName"`
		Type  string   `json:"mapType"`
		Next  []string `json:"mapNext"`
		Items []string `json:"mapItems"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	id := c.Query("targetMapId")

	pmap, _ := mapDao.Get(id)

	var list = []data{}
	row := data{}
	row.Id = pmap.Id
	row.Name = pmap.Name
	row.Type = pmap.Type
	row.Next = strings.Split(pmap.Next, ",")
	row.Items = strings.Split(pmap.Items, ",")
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func (*MapService) UpdateMapEndpoint(c *gin.Context) {

	type data struct {
		Id    string   `json:"mapId"`
		Name  string   `json:"mapName"`
		Type  string   `json:"mapType"`
		Next  []string `json:"mapNext"`
		Items []string `json:"mapItems"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	id := c.Query("targetPlayerId")

	items := c.Query("newMapItems")

	mapDao.Update(id, items)

	pmap, _ := mapDao.Get(id)

	var list = []data{}
	row := data{}
	row.Id = pmap.Id
	row.Name = pmap.Name
	row.Type = pmap.Type
	row.Next = strings.Split(pmap.Next, ",")
	row.Items = strings.Split(pmap.Items, ",")
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}
