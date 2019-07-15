package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/golang/glog"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var prodDBAddress, devDBAddress string

func queryDB(c echo.Context) error {
	var db *sql.DB
	var err error
	var res map[string]string

	res = make(map[string]string)

	env := c.Param("env")

	if env == "dev" {
		db, err = sql.Open("mysql", "dev:dev@172.20.0.191/dev?charset=utf8")
	} else {
		db, err = sql.Open("mysql", "prod:prod@172.20.0.192/prod?charset=utf8")
	}

	if err != nil {
		glog.Errorf("error: %s\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	rows, err := db.Query("SELECT * FROM userinfo")
	if err != nil {
		glog.Errorf("error: %s\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string

		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			glog.Errorf("error: %s\n", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		res["uid"] = strconv.Itoa(uid)
		res["username"] = username
		res["department"] = department
		res["created"] = created
	}

	db.Close()

	return c.JSON(http.StatusOK, res)
}

func main() {
	e := echo.New()

	// Get values from env
	prodDBAddress = os.Getenv("PROD_DB")
	devDBAddress = os.Getenv("DEV_DB")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/db/:env", queryDB)
	e.File("/", "public/index.html")
	e.Static("/", "assets")

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
