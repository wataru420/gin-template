package models

import (
	"database/sql"
	"log"
	"time"
)

type PlayerLog struct {
	Id    int
	PlayerId  string
	ApiPath  string
	ApiParam  string
	LogDateTime  string
}

type PlayerLogDao struct {
}

func (*PlayerLogDao) colums() string {
	return "playerLogId,playerId,apiPath,apiParam,logDatetime"
}

func (*PlayerLogDao) table() string {
	return "playerLog"
}

func (*PlayerLogDao) scan(rows sql.Rows) (PlayerLog, error) {
	rec := PlayerLog{}
	err := rows.Scan(&rec.Id, &rec.PlayerId, &rec.ApiPath, &rec.ApiParam, &rec.LogDateTime)
	return rec, err
}

func (dao *PlayerLogDao) GetList(id string) ([]PlayerLog, error) {
	var res = []PlayerLog{}
	rows, err := dbs.Query(`SELECT `+dao.colums()+`
							FROM `+dao.table()+` where playerId=? order by playerLogId desc limit 20`, id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return res, err
	}
	for rows.Next() {
		row,_ := dao.scan(*rows)
		res = append(res, row)
	}
	return res, err
}

func (dao *PlayerLogDao) Insert(playerId string, apiPath string, apiParam string) error {
	_, err := dbm.Exec(`INSERT INTO `+dao.table()+` (`+dao.colums()+`)
						VALUES (null,?,?,?,?)`,
		playerId,apiPath,apiParam,time.Now().Format("2006-01-02-15:04:05"))
	return err
}

