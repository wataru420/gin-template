package models
import (
	"log"
	"database/sql"
)

type Map struct {
	Id	  			string
	Name    		string
	Type	  		string
	Next		  	string
	Items		  	string
}

type MapDao struct {
}

func (*MapDao) colums() string {
	return "mapId,mapName,mapType,mapNext,mapItems"
}

func (*MapDao) table() string {
	return "map"
}

func (*MapDao) scan(rows sql.Rows) (Map, error) {
	rec := Map{}
	err := rows.Scan(&rec.Id,&rec.Name,&rec.Type, &rec.Next, &rec.Items)
	return rec, err
}

func (dao *MapDao) Get(id string) (Map, error) {
	rows, err := dbs.Query(`SELECT ` + dao.colums() + `
							FROM `+ dao.table() +` where mapId=?`, id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return Map{}, err
	}
	for rows.Next() {
		return dao.scan(*rows)
	}
	return Map{}, err
}

func (dao *MapDao) Update(id string, items string) error {
	query :="mapId='" + id + "'"
	if items != "" {
		query += ",mapItems='"+items + "'"
	}

	log.Println(query)
	_, err := dbm.Exec(`UPDATE `+ dao.table() +` SET
						` + query + `
						WHERE mapId=?`,
		id)
	return err
}

func (dao *MapDao) GetByItemId(id string) (Map, error) {
	id = "%" + id + "%"
	rows, err := dbs.Query(`SELECT `+dao.colums()+`
							FROM `+dao.table()+` where mapItems like ?`, id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return Map{}, err
	}
	for rows.Next() {
		return dao.scan(*rows)
	}
	return Map{}, err
}
