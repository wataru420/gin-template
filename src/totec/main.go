package main

import (
	"flag"
	"github.com/gin-gonic/gin"

	"totec/controllers"
	"totec/models"
)

func main() {

	var (
		apiBind                = flag.String("api-bind", ":8080", "Address to bind on")
		databaseMasterHost     = flag.String("database-master-host", "localhost", "Database master host")
		databaseMasterPort     = flag.Int("database-master-port", 3306, "Database master port")
		databaseMasterDbname   = flag.String("database-master-dbname", "totec", "Database master db name")
		databaseMasterUser     = flag.String("database-master-user", "root", "Database master username")
		databaseMasterPassword = flag.String("database-master-password", "", "Database master password")
		databaseSlaveHost      = flag.String("database-slave-host", "localhost", "Database slave host")
		databaseSlavePort      = flag.Int("database-slave-port", 3306, "Database slave port")
		databaseSlaveDbname    = flag.String("database-slave-dbname", "totec", "Database slave db name")
		databaseSlaveUser      = flag.String("database-slave-user", "root", "Database slave username")
		databaseSlavePassword  = flag.String("database-slave-password", "", "Database slave password")
		redisHost		       = flag.String("redis-host", "localhost", "Redis host")
		redisPort		       = flag.Int("reids-port", 6379, "Redis port")
	)

	flag.Parse()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	controllers.InitRooter(router)

	if err := models.Setup(&models.Config{
		DatabaseMasterHost:     *databaseMasterHost,
		DatabaseMasterPort:     *databaseMasterPort,
		DatabaseMasterDbname:   *databaseMasterDbname,
		DatabaseMasterUser:     *databaseMasterUser,
		DatabaseMasterPassword: *databaseMasterPassword,
		DatabaseSlaveHost:      *databaseSlaveHost,
		DatabaseSlavePort:      *databaseSlavePort,
		DatabaseSlaveDbname:    *databaseSlaveDbname,
		DatabaseSlaveUser:      *databaseSlaveUser,
		DatabaseSlavePassword:  *databaseSlavePassword,
		RedisHost:      		*redisHost,
		RedisPort:      		*redisPort,	}); err != nil {
		panic(err)
	}

	router.Run(*apiBind)
}
