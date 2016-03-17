package models

import (
	"database/sql"
	"log"
	"time"
)

type ItemLog struct {
	Id    int
	ItemId  string
	ApiPath  string
	ApiParam  string
	LogDateTime  string
}

type ItemLogDao struct {
}

func (*ItemLogDao) colums() string {
	return "itemLogId,itemId,apiPath,apiParam,logDatetime"
}

func (*ItemLogDao) table() string {
	return "itemLog"
}

func (*ItemLogDao) scan(rows sql.Rows) (ItemLog, error) {
	rec := ItemLog{}
	err := rows.Scan(&rec.Id, &rec.ItemId, &rec.ApiPath, &rec.ApiParam, &rec.LogDateTime)
	return rec, err
}

func (dao *ItemLogDao) GetList(id string) ([]ItemLog, error) {
	var res = []ItemLog{}
	rows, err := dbs.Query(`SELECT `+dao.colums()+`
							FROM `+dao.table()+` where itemId=? order by itemLogId desc limit 20`, id)
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

func (dao *ItemLogDao) Insert(playerId string, apiPath string, apiParam string) error {
	_, err := dbm.Exec(`INSERT INTO `+dao.table()+` (`+dao.colums()+`)
						VALUES (null,?,?,?,?)`,
		playerId,apiPath,apiParam,time.Now().Format("2006-01-02-15:04:05"))
	return err
}

