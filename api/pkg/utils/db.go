package utils

import (
	"database/sql"
	"log"
)

func CreateDatabaseIfNotExist(dbName string) {

	log.Printf("Creating database '%v' if not exists.", dbName)

	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/message_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS "+dbName)
	if err != nil {
		panic(err)
	}
	
	_, err = db.Exec("USE "+dbName)
	if err != nil {
		panic(err)
	}
	
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER NOT NULL AUTO_INCREMENT,
			body VARCHAR(280),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);`)

	if err != nil {
		panic(err)
	}
}