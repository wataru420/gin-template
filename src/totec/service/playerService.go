package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"totec/models"
)

type PlayerService struct{}

var playerDao = &models.PlayerDao{}

func (*PlayerService) ReadPlayerEndpoint(c *gin.Context) {

	type data struct {
		Id    string   `json:"playerId"`
		Name  string   `json:"playerName"`
		Hp    int      `json:"playerHp"`
		Mp    int      `json:"playerMp"`
		Exp   int      `json:"playerExp"`
		Atk   int      `json:"playerAtk"`
		Def   int      `json:"playerDef"`
		Int   int      `json:"playerInt"`
		Agi   int      `json:"playerAgi"`
		Items []string `json:"playerItems"`
		Map   string   `json:"playerMap"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	id := c.Query("targetPlayerId")

	player, _ := playerDao.Get(id)

	var list = []data{}
	row := data{}
	row.Id = player.Id
	row.Name = player.Name
	row.Hp = player.Hp
	row.Mp = player.Mp
	row.Exp = player.Exp
	row.Atk = player.Atk
	row.Def = player.Def
	row.Int = player.Int
	row.Agi = player.Agi
	row.Items = strings.Split(player.Items, ",")
	row.Map = player.Map
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func (*PlayerService) UpdatePlayerEndpoint(c *gin.Context) {

	type data struct {
		Id    string   `json:"playerId,omitempty"`
		Name  string   `json:"playerName,omitempty"`
		Hp    int      `json:"playerHp,omitempty"`
		Mp    int      `json:"playerMp,omitempty"`
		Exp   int      `json:"playerExp,omitempty"`
		Atk   int      `json:"playerAtk,omitempty"`
		Def   int      `json:"playerDef,omitempty"`
		Int   int      `json:"playerInt,omitempty"`
		Agi   int      `json:"playerAgi,omitempty"`
		Items []string `json:"playerItems,omitempty"`
		Map   string   `json:"playerMap,omitempty"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	id := c.Query("targetPlayerId")

	hp := c.Query("newPlayerHp")
	mp := c.Query("newPlayerMp")
	exp := c.Query("newPlayerExp")
	atk := c.Query("newPlayerAtk")
	def := c.Query("newPlayerDef")
	int := c.Query("newPlayerInt")
	agi := c.Query("newPlayerAgi")
	items := c.Query("newPlayerItems")
	playermap := c.Query("newPlayerMap")

	playerDao.Update(id, hp, mp, exp, atk, def, int, agi, items, playermap)

	param := playerLogParam{}
	iparam := itemLogParam{}
	param.Id = id
	if hp != "" {
		param.Hp, _ = strconv.Atoi(hp)
		iparam.Hp, _ = strconv.Atoi(hp)
	}
	if mp != "" {
		param.Mp, _ = strconv.Atoi(mp)
		iparam.Mp, _ = strconv.Atoi(mp)
	}
	if exp != "" {
		param.Exp, _ = strconv.Atoi(exp)
		iparam.Exp, _ = strconv.Atoi(exp)
	}
	if atk != "" {
		param.Atk, _ = strconv.Atoi(atk)
		iparam.Atk, _ = strconv.Atoi(atk)
	}
	if def != "" {
		param.Def, _ = strconv.Atoi(def)
		iparam.Def, _ = strconv.Atoi(def)
	}
	if int != "" {
		param.Int, _ = strconv.Atoi(int)
		iparam.Int, _ = strconv.Atoi(int)
	}
	if agi != "" {
		param.Agi, _ = strconv.Atoi(agi)
		iparam.Agi, _ = strconv.Atoi(agi)
	}
	if items != "" {
		param.Items = strings.Split(items, ",")
		iparam.Items = strings.Split(items, ",")
	}
	if playermap != "" {
		param.Map = playermap
		iparam.Map = playermap
	}
	bytes, _ := json.Marshal(param)
	playerLogDao.Insert(id, "updatePlayer", string(bytes))


	bytes, _ = json.Marshal(iparam)
	itemLogDao.Insert(id, "updatePlayer", string(bytes))

	player, _ := playerDao.Get(id)

	var list = []data{}
	row := data{}
	row.Id = player.Id
	row.Name = player.Name
	row.Hp = player.Hp
	row.Mp = player.Mp
	row.Exp = player.Exp
	row.Atk = player.Atk
	row.Def = player.Def
	row.Int = player.Int
	row.Agi = player.Agi
	row.Items = strings.Split(player.Items, ",")
	row.Map = player.Map
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func (*PlayerService) FindItemOwnerEndpoint(c *gin.Context) {

	id := c.Query("targetItemId")

	player, _ := playerDao.GetByItemId(id)

	if (player == models.Player{}) {

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

		pmap, _ := mapDao.GetByItemId(id)

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

	} else {
		returnPlayer(player, c)
	}

}

func (*PlayerService) SwitchItemOwnerEndpoint(c *gin.Context) {

	id := c.Query("targetItemId")
	owner := c.Query("newItemOwner")

	param := playerLogParam{}
	param.ItemId = id
	param.ItemOwner = owner

	bytes, _ := json.Marshal(param)
	playerLogDao.Insert(id, "switchItemOwner", string(bytes))


	player, _ := playerDao.GetByItemId(id)
	pmap, _ := mapDao.GetByItemId(id)
	if (pmap != models.Map{}) {
		items := strings.Split(pmap.Items, ",")
		items = remove(items, id)
		mapDao.Update(id, strings.Join(items, ","))
	}
	if (player != models.Player{}) {
		items := strings.Split(player.Items, ",")
		items = remove(items, id)
		playerDao.UpdateItems(id, strings.Join(items, ","))
	}

	iparam := itemLogParam{}
	iparam.ItemId = id
	iparam.ItemOwner = owner

	bytes, _ = json.Marshal(iparam)
	err := itemLogDao.Insert(id, "switchItemOwner", string(bytes))
	if err != nil {
		log.Println(err)
	}

	if owner == "none" {
		type data struct {
		}
		type res struct {
			Result bool   `json:"result"`
			Data   []data `json:"data"`
		}

		var list = []data{}
		response := res{true, list}
		log.Println(response)
		c.JSON(http.StatusOK, response)
	} else if strings.Index(owner, "Us") == 0 {
		log.Println("User:" + owner)
		player, _ := playerDao.Get(owner)
		itemsString := player.Items
		if player.Items != "" {
			items := strings.Split(player.Items, ",")
			items = append(items, id)
			itemsString = strings.Join(items, ",")
		}
		playerDao.UpdateItems(owner, itemsString)
		player.Items = itemsString
		returnPlayer(player, c)
	} else {
		log.Println("Map:" + owner)
		pmap, _ := mapDao.Get(owner)
		itemsString := pmap.Items
		if pmap.Items != "" {
			items := strings.Split(pmap.Items, ",")
			items = append(items, id)
			itemsString = strings.Join(items, ",")
		}
		mapDao.Update(owner, itemsString)
		pmap.Items = itemsString
		returnMap(pmap, c)
	}

}

func returnPlayer(player models.Player, c *gin.Context) {
	type data struct {
		Id    string   `json:"playerId"`
		Name  string   `json:"playerName"`
		Hp    int      `json:"playerHp"`
		Mp    int      `json:"playerMp"`
		Exp   int      `json:"playerExp"`
		Atk   int      `json:"playerAtk"`
		Def   int      `json:"playerDef"`
		Int   int      `json:"playerInt"`
		Agi   int      `json:"playerAgi"`
		Items []string `json:"playerItems"`
		Map   string   `json:"playerMap"`
	}
	type res struct {
		Result bool   `json:"result"`
		Data   []data `json:"data"`
	}

	var list = []data{}
	row := data{}
	row.Id = player.Id
	row.Name = player.Name
	row.Hp = player.Hp
	row.Mp = player.Mp
	row.Exp = player.Exp
	row.Atk = player.Atk
	row.Def = player.Def
	row.Int = player.Int
	row.Agi = player.Agi
	row.Items = strings.Split(player.Items, ",")
	row.Map = player.Map
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}

func remove(strings []string, search string) []string {
	result := []string{}
	for _, str := range strings {
		if str != search {
			result = append(result, str)
		}
	}
	return result
}
