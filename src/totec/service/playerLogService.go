package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"totec/models"
	"encoding/json"
)

type PlayerLogService struct{}

var playerLogDao = &models.PlayerLogDao{}

	type playerLogParam struct {
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
	}

func (*PlayerLogService) GetPlayerLogEndpoint(c *gin.Context) {
	id := c.Query("targetPlayerId")
	list, _ := playerLogDao.GetList(id)
	returnPlayerLog(list, c)
}

func returnPlayerLog(logList []models.PlayerLog, c *gin.Context) {

	type data struct {
		PlayerId    string `json:"playerId"`
		ApiPath     string `json:"apiPath"`
		ApiParam    playerLogParam `json:"apiParam"`
		LogDateTime string `json:"logDateTime"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	var list = []data{}
	for _, log := range logList {
		row := data{}
		row.PlayerId = log.PlayerId
		row.ApiPath = log.ApiPath
		json.Unmarshal([]byte(log.ApiParam),&row.ApiParam)
		row.LogDateTime = log.LogDateTime
		list = append(list, row)

	}

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}
