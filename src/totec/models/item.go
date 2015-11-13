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

func (*ItemDao) Get(id string) (Item, error) {
	res := Item{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "item:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT id,no,supplier,sold_quantity,sale_price,tags,image FROM items WHERE id=?`, id).Scan(
								&res.Id,&res.No,&res.Supplier,&res.SoldQuantity,&res.SalePrice,&res.Tags,&res.Image)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "item" + id,serialized)
	return res, err
}
