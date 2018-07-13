package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	configFile = "db_config.json"

	listPointStatment    = `SELECT * FROM point`
	insertPointStatement = `
		INSERT INTO point (
			bus_id, longitude, latitude, timestamp
		) VALUES (?, ?, ?, ?)`
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

	if err = prepareSQL(d); err != nil {
		return err
	}

	return nil
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

func prepareSQL(d *DBManager) (err error) {
	d.listPoint, err = d.db.Prepare(listPointStatment)
	if err != nil {
		return err
	}

	d.insertPoint, err = d.db.Prepare(insertPointStatement)
	if err != nil {
		return err
	}

	return nil
}
