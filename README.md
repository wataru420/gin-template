# gin-template

this is standard web app template made by GIN

## set up
install gb

```
$ go get github.com/constabulary/gb/...
```

restore dependencies from the manifest

```
$ gb vendor restore
```

build your application

```
$ gb build all
```

## run

just run your build file.

## option

you can change the connection of datasource using option

```
Usage:
  -api-bind string
    	Address to bind on (default ":8080")
  -database-master-dbname string
    	Database master db name (default "totec")
  -database-master-host string
    	Database master host (default "localhost")
  -database-master-password string
    	Database master password
  -database-master-port int
    	Database master port (default 3306)
  -database-master-user string
    	Database master username (default "root")
  -database-slave-dbname string
    	Database slave db name (default "totec")
  -database-slave-host string
    	Database slave host (default "localhost")
  -database-slave-password string
    	Database slave password
  -database-slave-port int
    	Database slave port (default 3306)
  -database-slave-user string
    	Database slave username (default "root")
  -redis-host string
    	Redis host (default "localhost")
  -reids-port int
    	Redis port (default 6379)
```