package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type User struct {
	Id	  			string
	No     			int
	PublicScore 	int
	Friends	  		string
	Image		  	string
}

type UserDao struct {
}

//func (*UserDao) GetList() ([]User, error) {
//	var res = []User{}
//	rows, err := dbs.Query(`SELECT id,name,password,type FROM users`)
//	if err != nil {
//		return res, err
//	}
//
//	for rows.Next() {
//		user := User{}
//		if err := rows.Scan(&user.Id,&user.Name,&user.Password,&user.Type); err != nil {
//			log.Fatal(err)
//			return res, err
//		}
//		res = append(res, user)
//	}
//	return res, err
//}

func (*UserDao) Get(id string) (User, error) {
	res := User{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "user:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT id,no,public_score,friends,image FROM users WHERE id=?`, id).Scan(&res.Id,&res.No,&res.PublicScore,&res.Friends,&res.Image)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "user:" + id,serialized)
	return res, err
}
