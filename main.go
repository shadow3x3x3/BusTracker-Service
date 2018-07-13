package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// AppTag can be use to debug tag
	AppTag = "shadow3x3x3/BusTracker"
)

var (
	mysqlDB DBManager
)

type trackerPoint struct {
	BusID     string  `json:"bus_id"`
	TimeStamp string  `json:"timestamp"`
	Longitude float64 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

func getTrackers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   "null", // TODO: Fetch data from db
	})
}

func postTrackers(c *gin.Context) {
	p := trackerPoint{}

	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}
	// TODO: Error handles

	fmt.Printf("Point: %+v\n", p)

	// TODO: Insert data into db
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func main() {
	if err := mysqlDB.Init(); err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	r.GET("/trackers", getTrackers)
	r.POST("/trackers", postTrackers)

	r.Run()
}
