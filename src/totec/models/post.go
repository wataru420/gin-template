package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"database/sql"
	"github.com/gin-gonic/gin"
	"strings"
	"strconv"

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









func (*PostDao) FindByParam(c *gin.Context, limit string) ([]Post, error) {
	var res = []Post{}

	var query = "SELECT " + allPostColums + " FROM posts WHERE 1=1"

	findByPostId := c.Query("findByPostId")
	if findByPostId != "" {
		query += " and id = '" + findByPostId + "'"
	}
/*
	findByItemSupplier := c.Query("findByItemSupplier")
	if findByItemSupplier != "" {
		query += " and supplier = '" + findByItemSupplier + "'"
	}

	findByItemSoldQuantityGTE := c.Query("findByItemSoldQuantityGTE")
	if findByItemSoldQuantityGTE != "" {
		query += " and sold_quantity >= " + findByItemSoldQuantityGTE
	}

	findByItemSoldQuantityLTE := c.Query("findByItemSoldQuantityLTE")
	if findByItemSoldQuantityLTE != "" {
		query += " and sold_quantity <= " + findByItemSoldQuantityLTE
	}

	findByItemSalePriceGTE := c.Query("findByItemSalePriceGTE")
	if findByItemSalePriceGTE != "" {
		query += " and sales_price >=" + findByItemSalePriceGTE
	}

	findByItemSalePriceLTE := c.Query("findByItemSalePriceLTE")
	if findByItemSalePriceLTE != "" {
		query += " and sales_price <=" + findByItemSalePriceLTE
	}

	findByItemTagsIncludeAll := c.Query("findByItemTagsIncludeAll")
	if findByItemTagsIncludeAll != "" {
		query += `
	     and id IN (
    	SELECT item_tags.tag
    	FROM item_tags
    	WHERE item_tags.id IN ("` + strings.Join(strings.Split(findByItemTagsIncludeAll,","),`","`) + `")
    	GROUP BY item_tags.tag
	    HAVING COUNT(item_tags.tag) >=
		` + strconv.Itoa(len(strings.Split(findByItemTagsIncludeAll,","))) + ")"
	}

	findByItemTagsIncludeAny := c.Query("findByItemTagsIncludeAny")
	if findByItemTagsIncludeAny != "" {
		query += `
	     and id NOT IN (
    	SELECT item_tags.tag
    	FROM item_tags
    	WHERE item_tags.id IN ("` + strings.Join(strings.Split(findByItemTagsIncludeAny,","),`","`) + `")
    	GROUP BY item_tags.tag
	    HAVING COUNT(item_tags.tag) >=
		` + strconv.Itoa(len(strings.Split(findByItemTagsIncludeAny,","))) + ")"
	}

	query += createScenario2PostQuery(c)
*/
	query += " limit " + limit

	rows, err := dbs.Query(query)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		row := Post{}
		if err := rows.Scan(&row.Id,&row.DateTime,&row.UserId,&row.ItemId,&row.ItemScore,&row.ItemState,&row.LikeUsers,&row.Tags); err != nil {
			return res, err
		}
		res = append(res, row)
	}
	return res, err
}

func createScenario2PostQuery(c *gin.Context) string {
	var scenario2 = false
	var query = " and id IN (SELECT item_id FROM posts WHERE 1=1 "
	findByPostId := c.Query("findByPostId")
	if findByPostId != "" {
		scenario2 = true
		query += ` and posts.id="` + findByPostId + `"`
	}

	findByPostDateTimeGTE := c.Query("findByPostDateTimeGTE")
	if findByPostDateTimeGTE != "" {
		scenario2 = true
		query += ` and posts.date_time>="` + findByPostDateTimeGTE + `"`
	}

	findByPostDateTimeLTE := c.Query("findByPostDateTimeLTE")
	if findByPostDateTimeLTE != "" {
		scenario2 = true
		query += ` and posts.date_time<="` + findByPostDateTimeLTE + `"`
	}

	findByPostItemId := c.Query("findByPostItemId")
	if findByPostItemId != "" {
		scenario2 = true
		query += ` and posts.item_id="` + findByPostItemId + `"`
	}

	findByMaxPostItemScoreGTE := c.Query("findByMaxPostItemScoreGTE")
	if findByMaxPostItemScoreGTE != "" {
		scenario2 = true
		query += ` and posts.item_score>="` + findByMaxPostItemScoreGTE +`"`
	}

	findByMaxPostItemScoreLTE := c.Query("findByMaxPostItemScoreLTE")
	if findByMaxPostItemScoreLTE != "" {
		scenario2 = true
		query += ` and posts.item_score<="` + findByMaxPostItemScoreLTE +`"`
	}

	findByPostItemState := c.Query("findByPostItemState")
	if findByPostItemState != "" {
		scenario2 = true
		query += ` and posts.item_state = "` + findByPostItemState + `"`
	}

	findByPostItemStateNotEQ := c.Query("findByPostItemStateNotEQ")
	if findByPostItemStateNotEQ != "" {
		scenario2 = true
		query += ` and posts.item_state != "` + findByPostItemStateNotEQ + `"`
	}

	findByPostLikeUsersIncludeUserIds := c.Query("findByPostLikeUsersIncludeUserIds")
	if findByPostLikeUsersIncludeUserIds != "" {
		scenario2 = true
		query += ` and id IN (
      SELECT post_likes.user_id
      FROM post_likes
      WHERE post_likes.id IN ("`+ strings.Join(strings.Split(findByPostLikeUsersIncludeUserIds,","),`","`)+`)
      GROUP BY post_likes.user_id
      HAVING COUNT(post_likes.user_id) >=
      ` + strconv.Itoa(len(strings.Split(findByPostLikeUsersIncludeUserIds,","))) + `")`
	}

	findByPostLikeUsersNotIncludeUserIds:= c.Query("findByPostLikeUsersNotIncludeUserIds")
	if findByPostLikeUsersNotIncludeUserIds != "" {
		scenario2 = true
		query += ` and id IN (
      SELECT post_likes.user_id
      FROM post_likes
      WHERE post_likes.id IN ("`+ strings.Join(strings.Split(findByPostLikeUsersNotIncludeUserIds,","),`","`)+`)
      GROUP BY post_likes.user_id
      HAVING COUNT(post_likes.user_id) >=
      ` + strconv.Itoa(len(strings.Split(findByPostLikeUsersNotIncludeUserIds,","))) + `")`
	}
	if scenario2 {
		return query + ")"
	} else {
		return ""
	}
}
