package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// AppTag can be use to debug tag
	AppTag = "shadow3x3x3/BusTracker"
)

var (
	mysqlDB DBManager
)

// TrackerPoint is a basic unit
type TrackerPoint struct {
	BusID     string  `json:"bus_id"`
	TimeStamp string  `json:"timestamp"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func getTrackers(c *gin.Context) {
	points, err := mysqlDB.ListPoints()

	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   points,
	})
}

func postTrackers(c *gin.Context) {
	p := TrackerPoint{}

	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "json decoding : " + err.Error(),
			"status": http.StatusBadRequest,
		})
		return
	}
	// TODO: Error handles

	fmt.Printf("Point: %+v\n", p)

	_, err := mysqlDB.insertPoint.Query(p.BusID, p.Latitude, p.Longitude, p.TimeStamp)

	if err != nil {
		fmt.Println(err)
		return
	}

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
