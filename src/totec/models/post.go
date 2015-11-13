package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"database/sql"
)

type Post struct {
	Id	  			string
	DateTime		int
	UserId  		string
	ItemId  		string
	ItemScore	 	int
	ItemState  		string
	LikeUsers  		string
	Tags		  	string
}

type PostDao struct {
}

var allPostColums = "id,date_time,user_id,item_id,score,state,like_users,tags"

func (*PostDao) Get(id string) (Post, error) {
	res := Post{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "post:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT ` + allPostColums + ` FROM posts WHERE id=?`, id).Scan(
						&res.Id,&res.DateTime,&res.UserId,&res.ItemId,&res.ItemScore,&res.ItemState,&res.LikeUsers,&res.Tags)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "item" + id,serialized)
	return res, err
}

func (*PostDao) FindByPostUserId(id string, limit int) ([]Post, error) {
	rows, err := dbs.Query(`SELECT ` + allPostColums + ` FROM posts where user_id=? limit ?`, id, limit)
	if err != nil {
		return []Post{}, err
	}
	return postRowScan(rows)
}


func (*PostDao) FindByPostItemId(id string, limit int) ([]Post, error) {
	rows, err := dbs.Query(`SELECT ` + allPostColums + ` FROM posts where item_id=? limit ?`, id, limit)
	if err != nil {
		return []Post{}, err
	}
	return postRowScan(rows)
}

func postRowScan(rows *sql.Rows) ([]Post,error) {
	var res = []Post{}
	for rows.Next() {
		row := Post{}
		if err := rows.Scan(&row.Id,&row.DateTime,&row.UserId,&row.ItemId,&row.ItemScore,&row.ItemState,&row.LikeUsers,&row.Tags); err != nil {
			return res, err
		}
		res = append(res, row)
	}
	return res, nil
}
