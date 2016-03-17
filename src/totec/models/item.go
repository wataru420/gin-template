package models
import (
	"log"
	"database/sql"
)

type Item struct {
	Id	  			string
	Name    		string
	Type	  		string
	Value 	int
	EffectTarget		  	string
	EffectValue 	int
}

type ItemDao struct {
}

func (*ItemDao) colums() string {
	return "itemId,itemName,itemType,itemValue,itemEffectTarget,itemEffectValue"
}

func (*ItemDao) table() string {
	return "item"
}

func (*ItemDao) scan(rows sql.Rows) (Item, error) {
	rec := Item{}
	err := rows.Scan(&rec.Id,&rec.Name,&rec.Type, &rec.Value, &rec.EffectTarget,&rec.EffectValue)
	return rec, err
}

func (dao *ItemDao) Get(id string) (Item, error) {
	rows, err := dbs.Query(`SELECT ` + dao.colums() + `
							FROM `+ dao.table() +` where itemId=?`, id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return Item{}, err
	}
	for rows.Next() {
		return dao.scan(*rows)
	}
	return Item{}, err
}

func (dao *ItemDao) Update(id string, value string) error {
	query :="itemId='" + id + "'"
	if value != "" {
		query += ",itemValue="+value
	}

	log.Println(query)
	_, err := dbm.Exec(`UPDATE `+ dao.table() +` SET
						` + query + `
						WHERE itemId=?`,
		id)
	return err
}
