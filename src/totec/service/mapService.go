package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"totec/models"
	"encoding/json"
)

type MapService struct{}

var mapDao = &models.MapDao{}

func (*MapService) ReadMapEndpoint(c *gin.Context) {
	id := c.Query("targetMapId")
	pmap, _ := mapDao.Get(id)
	returnMap(pmap, c)
}

func (*MapService) UpdateMapEndpoint(c *gin.Context) {
	id := c.Query("targetMapId")
	items := c.Query("newMapItems")
	mapDao.Update(id, items)

	iparam := itemLogParam{}
	iparam.MapId = id
	iparam.MapItems = strings.Split(items,",")

	bytes, _ := json.Marshal(iparam)
	err := itemLogDao.Insert(id, "switchItemOwner", string(bytes))
	if err != nil {
		log.Println(err)
	}


	pmap, _ := mapDao.Get(id)

	returnMap(pmap, c)
}

func returnMap(pmap models.Map, c *gin.Context) {

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
