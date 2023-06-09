package server

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	onceDb sync.Once
	db     *sql.DB
)

func NewDb() *sql.DB {
	onceDb.Do(func() {
		connectionString := "goaplication:goIsKing@tcp(127.0.0.1:3306)/humservice"
		var err error
		db, err = sql.Open("mysql", connectionString)

		if err != nil {
			panic(err)
		}

		// defer db.Close()

		fmt.Println("DB Connected")

		initTables(db)
	})

	return db
}

func initTables(db *sql.DB) {
	stationQuery := `CREATE TABLE IF NOT EXISTS Stations (
		uuid VARCHAR(255), 
		url VARCHAR(255),
		created DATETIME,
		PRIMARY KEY (uuid)
		)`

	initTable(db, stationQuery)

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

func initTable(db *sql.DB, query string) {
	_, err := db.Exec(query)

	if err != nil {
		panic(err)
	}
}
