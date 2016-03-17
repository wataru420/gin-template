package main

import (
	"fmt"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/BurntSushi/toml"

	"totec/controllers"
	"totec/models"
)

type Config struct {
	Api         ApiConfig         `toml:"api"`
	Database    DatabaseConfig    `toml:"database"`
	Redis       RedisConfig       `toml:"redis"`
}

type DatabaseConfig struct {
	Master DatabaseServerConfig `toml:"master"`
	Slave  DatabaseServerConfig `toml:"slave"`
}

type DatabaseServerConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Dbname   string `toml:"dbname"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
}

type ApiConfig struct {
	Port int `toml:port`
}

var Conf Config

func main() {

	var (
		confFile = flag.String("conf", "config.tml", "config toml file")
	)

	flag.Parse()

	_, err := toml.DecodeFile(*confFile, &Conf)
	if err != nil {
		panic(err)
	}

	flag.Parse()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	controllers.InitRooter(router)

	if err := models.Setup(&models.Config{
		DatabaseMasterHost:     Conf.Database.Master.Host,
		DatabaseMasterPort:     Conf.Database.Master.Port,
		DatabaseMasterDbname:   Conf.Database.Master.Dbname,
		DatabaseMasterUser:     Conf.Database.Master.User,
		DatabaseMasterPassword: Conf.Database.Master.Password,
		DatabaseSlaveHost:      Conf.Database.Slave.Host,
		DatabaseSlavePort:      Conf.Database.Slave.Port,
		DatabaseSlaveDbname:    Conf.Database.Slave.Dbname,
		DatabaseSlaveUser:      Conf.Database.Slave.User,
		DatabaseSlavePassword:  Conf.Database.Slave.Password,
		RedisHost:      		Conf.Redis.Host,
		RedisPort:      		Conf.Redis.Port,	}); err != nil {
		panic(err)
	}

	router.Run(fmt.Sprintf(":%d", Conf.Api.Port))
}
