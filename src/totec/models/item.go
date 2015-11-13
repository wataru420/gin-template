package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type Item struct {
	Id	  			string
	No     			int
	Supplier  		string
	SoldQuantity 	int
	SalePrice	 	int
	Tags	  		string
	Image		  	string
}

type ItemDao struct {
}

var allColums = "id,no,supplier,solid_quantity,sale_price,tags,image"
func (*ItemDao) Get(id string) (Item, error) {
	res := Item{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "item:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT ` + allColums + ` FROM items WHERE id=?`, id).Scan(
								&res.Id,&res.No,&res.Supplier,&res.SoldQuantity,&res.SalePrice,&res.Tags,&res.Image)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "item" + id,serialized)
	return res, err
}

func (*ItemDao) FindByPostUserId(id string, limit int) ([]Item, error) {
	var res = []Item{}
	rows, err := dbs.Query(`SELECT ` + allColums + ` FROM items where id=? limit ?`, id, limit)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		item := Item{}
		if err := rows.Scan(&item.Id,&item.No,&item.Supplier,&item.SoldQuantity,&item.SalePrice,&item.Tags,&item.Image); err != nil {
			return res, err
		}
		res = append(res, item)
	}
	return res, err
}
