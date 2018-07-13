package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	appTag     = "BusTracker"
	configFile = "db_config.json"
)

type databaseConfig struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

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

func readDatabaseConfig(c *databaseConfig) error {
	file, err := os.Open(configFile)

	defer file.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&c); err != nil {
		return err
	}

	return nil
}

func initDatabase() error {
	config := databaseConfig{}

	if err := readDatabaseConfig(&config); err != nil {
		return err
	}

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.IP, config.Port, config.Database)

	db, err := sql.Open("mysql", dbURL)

	if err != nil {
		return err
	}

	// Simple test db connection
	if _, err := db.Query("select 1"); err != nil {
		return err
	}

	fmt.Printf("[%s] Database is OK!\n\n", appTag)

	return nil
}

func main() {
	if err := initDatabase(); err != nil {
		fmt.Println(err)
		return
	}

	r := gin.Default()

	r.GET("/trackers", getTrackers)
	r.POST("/trackers", postTrackers)

	r.Run()
}
