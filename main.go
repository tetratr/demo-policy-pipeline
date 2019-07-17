package main

import (
	"flag"
	"fmt"
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

func status(c echo.Context) error {
	var res map[string]string
	res = make(map[string]string)

	res["environment"] = os.Getenv("ENV_LIFECYCLE")
	res["commit_id"] = os.Getenv("COMMIT_ID")

	return c.JSON(http.StatusOK, res)
}

func queryDB(c echo.Context) error {
	var db *sql.DB
	var err error
	var res map[string]string

	res = make(map[string]string)

	env := c.Param("env")

	if env == "dev" {
		db, err = sql.Open("mysql", fmt.Sprintf("dev:dev@tcp(%s:3306)/dev?charset=utf8&timeout=5s", devDBAddress))
	} else {
		db, err = sql.Open("mysql", fmt.Sprintf("prod:prod@tcp(%s:3306)/prod?charset=utf8&timeout=5s", prodDBAddress))
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
		var uid sql.NullInt64
		var username sql.NullString
		var department sql.NullString
		var created sql.NullString

		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			glog.Errorf("error: %s\n", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		res["uid"] = strconv.Itoa(int(uid.Int64))
		res["username"] = username.String
		res["department"] = department.String
		res["created"] = created.String
	}

	db.Close()

	return c.JSON(http.StatusOK, res)
}

func main() {
	flag.Parse()

	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("2")

	e := echo.New()

	// Get values from env
	prodDBAddress = os.Getenv("PROD_DB")
	devDBAddress = os.Getenv("DEV_DB")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/db/:env", queryDB)
	e.GET("/status", status)
	e.Static("/", "web")
	e.File("/", "web/index.html")

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}
