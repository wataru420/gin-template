package models
import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"strconv"
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

var allItemColums = "items.id,items.no,items.supplier,items.solid_quantity,items.sale_price,items.tags,items.image "

func (*ItemDao) Get(id string) (Item, error) {
	res := Item{Id:id}
	con := redisPool.Get()
	defer con.Close()

	s, err := redis.Bytes(con.Do("GET", "item:" + id))
	if s != nil {
		json.Unmarshal(s, &res)
		return res, err
	}

	err = dbs.QueryRow(`SELECT ` + allItemColums + ` FROM items WHERE id=?`, id).Scan(
								&res.Id,&res.No,&res.Supplier,&res.SoldQuantity,&res.SalePrice,&res.Tags,&res.Image)
	serialized, _ := json.Marshal(res)
	con.Do("SET", "item" + id,serialized)
	return res, err
}

func (*ItemDao) FindByPostUserId(id string, limit int) ([]Item, error) {
	var res = []Item{}
	rows, err := dbs.Query(`SELECT ` + allItemColums + `FROM items INNER JOIN posts ON items.id = posts.item_id where posts.user_id=? limit ?`, id, limit)
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

func (*ItemDao) FindByParam(c *gin.Context, limit string) ([]Item, error) {
	var res = []Item{}

	var query = "SELECT " + allItemColums + " FROM items WHERE 1=1"

	findByItemId := c.Query("findByItemId")
	if findByItemId != "" {
		query += " and id = '" + findByItemId + "'"
	}

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

	query += createScenario2ItemQuery(c)

	query += " limit " + limit

	log.Printf(query)
	rows, err := dbs.Query(query)
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

func createScenario2ItemQuery(c *gin.Context) string {
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
