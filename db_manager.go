package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbConfigFile = "db_config.json"

	listPointStatement = `
		SELECT bus_id, longitude, latitude, timestamp
		FROM point`

	insertPointStatement = `
		INSERT INTO point (
			bus_id, longitude, latitude, timestamp
		)
		VALUES (?, ?, ?, ?)`
)

type dbConfig struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// DBManager is responsible for transactions of database
type DBManager struct {
	db          *sql.DB
	listPoint   *sql.Stmt
	insertPoint *sql.Stmt
}

// Init method can initialize DBManager from config file.
func (d *DBManager) Init() (err error) {
	if d.db, err = initDatabase(); err != nil {
		return err
	}

	return prepareSQL(d)
}

// ListPoints can query all points from db and make them to TrackerPoint struct
func (d *DBManager) ListPoints() ([]*TrackerPoint, error) {
	points := []*TrackerPoint{}

	rows, err := d.listPoint.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		busID := ""
		longitude := float64(0)
		latitude := float64(0)
		timestamp := ""

		err = rows.Scan(&busID, &longitude, &latitude, &timestamp)

		if err != nil {
			return nil, err
		}

		points = append(points, &TrackerPoint{
			BusID:     busID,
			TimeStamp: timestamp,
			Longitude: longitude,
			Latitude:  latitude,
		})

	}

	return points, nil
}

func initDatabase() (*sql.DB, error) {
	config := dbConfig{}

	if err := readDatabaseConfig(&config); err != nil {
		return nil, err
	}

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.IP, config.Port, config.Database)

	db, err := sql.Open("mysql", dbURL)

	if err != nil {
		return nil, err
	}

	// Simple test db connection
	if _, err := db.Query("select 1"); err != nil {
		return nil, err
	}

	fmt.Printf("[%s] Database is OK!\n\n", AppTag)

	return db, nil
}

func readDatabaseConfig(c *dbConfig) error {
	file, err := os.Open(dbConfigFile)

	defer file.Close()

	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	return decoder.Decode(&c)
}

func prepareSQL(d *DBManager) (err error) {
	d.listPoint, err = d.db.Prepare(listPointStatement)
	if err != nil {
		return err
	}

	d.insertPoint, err = d.db.Prepare(insertPointStatement)
	if err != nil {
		return err
	}

	return nil
}
