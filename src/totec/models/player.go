package models

import (
	"database/sql"
	"log"
)

type Player struct {
	Id    string
	Name  string
	Hp    int
	Mp    int
	Exp   int
	Atk   int
	Def   int
	Int   int
	Agi   int
	Items string
	Map   string
}

type PlayerDao struct {
}

func (*PlayerDao) colums() string {
	return "playerId,playerName,playerHp,playerMp,playerExp,playerAtk,playerDef,playerInt,playerAgi,playerItems,playerMap"
}

func (*PlayerDao) table() string {
	return "player"
}

func (*PlayerDao) scan(rows sql.Rows) (Player, error) {
	rec := Player{}
	err := rows.Scan(&rec.Id, &rec.Name, &rec.Hp, &rec.Mp, &rec.Exp, &rec.Atk, &rec.Def, &rec.Int, &rec.Agi, &rec.Items, &rec.Map)
	return rec, err
}

func (dao *PlayerDao) Get(id string) (Player, error) {
	rows, err := dbs.Query(`SELECT `+dao.colums()+`
							FROM `+dao.table()+` where playerId=?`, id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return Player{}, err
	}
	for rows.Next() {
		return dao.scan(*rows)
	}
	return Player{}, err
}

func (dao *PlayerDao) GetByItemId(id string) (Player, error) {
	id = "%" + id + "%"
	rows, err := dbs.Query(`SELECT `+dao.colums()+`
							FROM `+dao.table()+` where playerItems like ?`, id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return Player{}, err
	}
	for rows.Next() {
		return dao.scan(*rows)
	}
	return Player{}, err
}

func (dao *PlayerDao) Update(id string, hp string, mp string, exp string, atk string, def string, int string, agi string, items string, playermap string) error {
	query := "playerId='" + id + "'"
	if hp != "" {
		query += ",playerHp=" + hp
	}
	if mp != "" {
		query += ",playerMp=" + mp
	}
	if exp != "" {
		query += ",playerExp=" + exp
	}
	if atk != "" {
		query += ",playerAtk=" + atk
	}
	if def != "" {
		query += ",playerDef=" + def
	}
	if int != "" {
		query += ",playerInt=" + int
	}
	if agi != "" {
		query += ",playerAgi=" + agi
	}
	if items != "" {
		query += ",playerItems='" + items + "'"
	}
	if playermap != "" {
		query += ",playerMap='" + playermap + "'"
	}

	log.Println(query)
	_, err := dbm.Exec(`UPDATE `+dao.table()+` SET
						`+query+`
						WHERE playerId=?`,
		id)
	return err
}

func (dao *PlayerDao) UpdateItems(id string, items string) error {
	query :=  "playerItems='" + items + "'"

	log.Println(query)
	_, err := dbm.Exec(`UPDATE `+dao.table()+` SET
						`+query+`
						WHERE playerId=?`,
		id)
	return err
}
