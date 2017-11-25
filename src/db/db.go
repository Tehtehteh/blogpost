package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB


func InitializeDB(connstring string) (err error) {
	DB, err = sql.Open("mysql", connstring)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func FetchOne (query string, channel chan *sql.Row, args ...interface{}){
	err := DB.Ping()
	if err != nil {
		log.Panicf("Error pinging database: %s\n", err)
	}
	res, err := DB.Prepare(query)

	if err != nil {
		log.Panicf("Error executing SQL query: %s", err)
	}

	defer res.Close()
	defer close(channel)

	channel <- res.QueryRow(args...)
}

func FetchMany(query string, channel chan *sql.Rows, args ...interface{}){
	err := DB.Ping()
	if err != nil {
		log.Panicf("[Mysql] Erorr: %s", err)
	}
	res, err := DB.Prepare(query)

	var rows *sql.Rows

	if err != nil {
		log.Panicf("[Mysql] Erorr: %s", err)
	}
	if len(args) != 0 {  //todo: Check correct way to handle empty args
		rows, err = res.Query(args)
	} else {
		rows, err = res.Query()
	}
	if err != nil {
		log.Panicf("[Mysql] Erorr: %s", err)
	}
	defer res.Close()
	defer close(channel)

	channel <- rows
}
