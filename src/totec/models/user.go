package models
import (
	"log"
	"encoding/json"
	"strconv"
	"github.com/garyburd/redigo/redis"
)

type User struct {
	Id     		int
	Name	  	string
	Password  	string
	Type	  	int
}

type UserDao struct {
}

func (*UserDao) GetList() ([]User, error) {
	var res = []User{}
	rows, err := dbs.Query(`SELECT id,name,password,type FROM users`)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Id,&user.Name,&user.Password,&user.Type); err != nil {
			log.Fatal(err)
			return res, err
		}
		res = append(res, user)
	}
	return res, err
}

func (*UserDao) Get(id int) (User, error) {
	res := User{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", strconv.Itoa(id)))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT name,password,type FROM users WHERE id=?`, id).Scan(&res.Name,&res.Password,&res.Type)
	serialized, _ := json.Marshal(res)
	con.Do("SET", strconv.Itoa(id),serialized)
	return res, err
}
