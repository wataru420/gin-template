package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
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

func (*PostDao) Get(id string) (Post, error) {
	res := Post{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "post:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT id,date_time,user_id,item_id,score,state,like_users,tags FROM items WHERE id=?`, id).Scan(
						&res.Id,&res.DateTime,&res.UserId,&res.ItemId,&res.ItemScore,&res.ItemState,&res.LikeUsers,&res.Tags)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "item" + id,serialized)
	return res, err
}
