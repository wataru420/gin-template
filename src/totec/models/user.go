package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
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

var allUserColums = "users.id,users.no,users.public_score,users.friends,users.image"

func (*UserDao) Get(id string) (User, error) {
	res := User{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "user:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT ` + allUserColums + ` FROM users WHERE id=?`, id).Scan(&res.Id,&res.No,&res.PublicScore,&res.Friends,&res.Image)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "user:" + id,serialized)
	return res, err
}

func (*UserDao) FindByPostItemId(id string, limit int) ([]User, error) {
	var res = []User{}
	rows, err := dbs.Query(`SELECT ` + allUserColums + `FROM users INNER JOIN posts ON user.id = posts.user_id where posts.item_id=? limit ?`, id, limit)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		row := User{}
		if err := rows.Scan(&row.Id,&row.No,&row.PublicScore,&row.Friends,&row.Image); err != nil {
			return res, err
		}
		res = append(res, row)
	}
	return res, err
}

func (*UserDao) FindByParam(c *gin.Context, limit string) ([]User, error) {
	var res = []User{}

	var query = "SELECT " + allUserColums + " FROM users WHERE "

	findByUserId := c.Query("findByUserId")
	if findByUserId != "" {
		query += " id = '" + findByUserId + "'"
	}

	findByUserPublicScoreGTE := c.Query("findByUserPublicScoreGTE")
	if findByUserPublicScoreGTE != "" {
		query += " public_score >= " + findByUserPublicScoreGTE
	}

	findByUserPublicScoreLTE := c.Query("findByUserPublicScoreLTE")
	if findByUserPublicScoreLTE != "" {
		query += " public_score <= " + findByUserPublicScoreLTE
	}

	findByUserFriendsNumberGTE := c.Query("findByUserFriendsNumberGTE")
	if findByUserFriendsNumberGTE != "" {
		query += " friend_num >=" + findByUserFriendsNumberGTE
	}

	findByUserFriendsNumberLTE := c.Query("findByUserFriendsNumberLTE")
	if findByUserFriendsNumberLTE != "" {
		query += " friend_num <=" + findByUserFriendsNumberLTE
	}

	query += " limit " + limit

	log.Printf(query)
	rows, err := dbs.Query(query)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		row := User{}
		if err := rows.Scan(&row.Id,&row.No,&row.PublicScore,&row.Friends,&row.Image); err != nil {
			return res, err
		}
		res = append(res, row)
	}
	return res, err
}
