package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Config struct {
	DatabaseMasterHost     string
	DatabaseMasterPort     int
	DatabaseMasterDbname   string
	DatabaseMasterUser     string
	DatabaseMasterPassword string
	DatabaseSlaveHost      string
	DatabaseSlavePort      int
	DatabaseSlaveDbname    string
	DatabaseSlaveUser      string
	DatabaseSlavePassword  string
	RedisHost			   string
	RedisPort			   int
}

var (
	dbm *sql.DB
	dbs *sql.DB
	redisPool *redis.Pool
)

func Setup(config *Config) error {
	connectString := "%s:%s@(%s:%d)/%s?timeout=%ds"

	dataSourceMaster := fmt.Sprintf(connectString,
		config.DatabaseMasterUser,
		config.DatabaseMasterPassword,
		config.DatabaseMasterHost,
		config.DatabaseMasterPort,
		config.DatabaseMasterDbname,
		3000,
	)

	dataSourceSlave := fmt.Sprintf(connectString,
		config.DatabaseSlaveUser,
		config.DatabaseSlavePassword,
		config.DatabaseSlaveHost,
		config.DatabaseSlavePort,
		config.DatabaseSlaveDbname,
		3000,
	)

	{
		var err error
		dbm, err = sql.Open("mysql", dataSourceMaster)
		if err != nil {
			return err
		}
		var res string
		if err := dbm.QueryRow("SELECT 1").Scan(&res); err != nil {
			panic(err)
		}
	}
	{
		var err error
		dbs, err = sql.Open("mysql", dataSourceSlave)
		if err != nil {
			return err
		}
		var res string
		if err := dbs.QueryRow("SELECT 1").Scan(&res); err != nil {
			panic(err)
		}
	}

	redisConnectString := "%s:%d"
	{
		redisPool = &redis.Pool{
			MaxIdle: 10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", fmt.Sprintf(redisConnectString,config.RedisHost,config.RedisPort))
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}
	return nil
}
