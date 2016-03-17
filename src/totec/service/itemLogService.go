package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"totec/models"
	"encoding/json"
)

type ItemLogService struct{}

var itemLogDao = &models.ItemLogDao{}

	type itemLogParam struct {
		Id    string   `json:"playerId,omitempty"`
		Hp    int      `json:"newPlayerHp,omitempty"`
		Mp    int      `json:"newPlayerMp,omitempty"`
		Exp   int      `json:"newPlayerExp,omitempty"`
		Atk   int      `json:"newPlayerAtk,omitempty"`
		Def   int      `json:"newPlayerDef,omitempty"`
		Int   int      `json:"newPlayerInt,omitempty"`
		Agi   int      `json:"newPlayerAgi,omitempty"`
		Items []string `json:"newPlayerItems,omitempty"`
		Map   string   `json:"newPlayerMap,omitempty"`


		ItemId    string   `json:"itemId,omitempty"`
		ItemOwner    string   `json:"newItemOwner,omitempty"`

		MapId	string   `json:"targetMapId,omitempty"`
		MapItems []string   `json:"newMapItems,omitempty"`
	}

func (*ItemLogService) GetItemLogEndpoint(c *gin.Context) {
	id := c.Query("targetItemId")
	list, _ := itemLogDao.GetList(id)
	returnItemLog(list, c)
}

func returnItemLog(logList []models.ItemLog, c *gin.Context) {

	type data struct {
		ItemId    string `json:"itemId"`
		ApiPath     string `json:"apiPath"`
		ApiParam    itemLogParam `json:"apiParam"`
		LogDateTime string `json:"logDateTime"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	var list = []data{}
	for _, log := range logList {
		row := data{}
		row.ItemId = log.ItemId
		row.ApiPath = log.ApiPath
		json.Unmarshal([]byte(log.ApiParam),&row.ApiParam)
		row.LogDateTime = log.LogDateTime
		list = append(list, row)

	}

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}
