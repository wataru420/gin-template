package service
import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"totec/models"
	"strings"
)


type PlayerService struct {}

var playerDao = &models.PlayerDao{}

func (*PlayerService) ReadPlayerEndpoint(c *gin.Context)  {

	type data struct {
		Id string	`json:"playerId"`
		Name string	`json:"playerName"`
		Hp int	`json:"playerHp"`
		Mp int	`json:"playerMp"`
		Exp int	`json:"playerExp"`
		Atk int	`json:"playerAtk"`
		Def int	`json:"playerDef"`
		Int int	`json:"playerInt"`
		Agi int	`json:"playerAgi"`
		Items []string	`json:"playerItems"`
		Map string	`json:"playerMap"`
	}
	type res struct {
		Result bool        `json:"result"`
		Data   []data `json:"data"`
	}


	id := c.Query("targetPlayerId")

	player,_ := playerDao.Get(id)

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
	row.Items = strings.Split(player.Items,",")
	row.Map = player.Map
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}


func (*PlayerService) UpdatePlayerEndpoint(c *gin.Context)  {

	type data struct {
		Id string	`json:"playerId"`
		Name string	`json:"playerName"`
		Hp int	`json:"playerHp"`
		Mp int	`json:"playerMp"`
		Exp int	`json:"playerExp"`
		Atk int	`json:"playerAtk"`
		Def int	`json:"playerDef"`
		Int int	`json:"playerInt"`
		Agi int	`json:"playerAgi"`
		Items []string	`json:"playerItems"`
		Map string	`json:"playerMap"`
	}
	type res struct {
		Result bool        `json:"result"`
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

	playerDao.Update(id,hp,mp,exp,atk,def,int,agi,items,playermap)

	player,_ := playerDao.Get(id)

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
	row.Items = strings.Split(player.Items,",")
	row.Map = player.Map
	list = append(list, row)

	response := res{true, list}
	log.Println(response)
	c.JSON(http.StatusOK, response)
}
