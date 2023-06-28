package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	onceDb sync.Once
	db     *sql.DB
)

// Singleton for DB Connection
func NewDb() *sql.DB {
	var err error
	db, err = sql.Open("mysql", os.Getenv("MYSQL_STRING"))
	onceDb.Do(func() {

		if err != nil {
			panic(err)
		}

		fmt.Println("DB Connected")

		// init Tables necessary for the Application
		initTables(db)
	})

	return db
}

// Inits Tables for the Application
func initTables(db *sql.DB) {
	// Table for all the Stations
	stationQuery := `CREATE TABLE IF NOT EXISTS Stations (
		uuid VARCHAR(255), 
		url VARCHAR(255),
		created DATETIME,
		PRIMARY KEY (uuid)
		)`

	initTable(db, stationQuery)

	// Table for all the data received from the stations
	dataQuery := `CREATE TABLE IF NOT EXISTS Data (
		id INT NOT NULL AUTO_INCREMENT, 
		hum DOUBLE(5,2),
		temp DOUBLE(5,2),
		time DATETIME,
		station VARCHAR(255),
		PRIMARY KEY (id),
		FOREIGN KEY (station) REFERENCES Stations(uuid)
		)`

	initTable(db, dataQuery)
}

// executes the given Create Statement at the DB
func initTable(db *sql.DB, query string) {
	_, err := db.Exec(query)

	if err != nil {
		panic(err)
	}
}
