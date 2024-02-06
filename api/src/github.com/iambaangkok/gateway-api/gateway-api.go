package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// "net/http"
	"github.com/go-martini/martini"
	// "github.com/martini-contrib/encoder"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func queryRandomMessage(res http.ResponseWriter, req *http.Request, params martini.Params) []byte {

	var msg Message

	err := db.QueryRow(
		`SELECT *
		FROM messages
		ORDER BY RAND()
		LIMIT 1`).
		Scan(&msg.ID, &msg.Body, &msg.CreatedAt)
  
	if err != nil && err != sql.ErrNoRows {
	  log.Println(err)
	}
	log.Print(msg.toString())

	return msgToJson(msg)
}

func queryMessageById(id int64) Message {

	var msg Message

	err := db.QueryRow(
		`SELECT *
		FROM messages m
		WHERE m.id = ?
		LIMIT 1`, id).
		Scan(&msg.ID, &msg.Body, &msg.CreatedAt)
  
	if err != nil && err != sql.ErrNoRows {
	  log.Println(err)
	}
	log.Print(msg.toString())

	return msg
}

func msgToJson(msg Message) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return b
}

func insertMessage(res http.ResponseWriter, req *http.Request, params martini.Params) []byte {
	decoder := json.NewDecoder(req.Body)

	var dto MessageInputDTO
	err := decoder.Decode(&dto)
	if err != nil {
		panic(err)
	}

	insertStatement := fmt.Sprintf(
		`INSERT INTO messages (body)
		VALUES ("%v")`, dto.Body)

	result, err := db.ExecContext(context.Background(), insertStatement)
	if err != nil {
		log.Println(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	msg := queryMessageById(id)

	return msgToJson(msg)
}

func handleRequests() {
	m := martini.Classic()

	m.Get("/messages/random", queryRandomMessage)
	m.Post("/messages/add", insertMessage)
	m.Run()
}

type MessageInputDTO struct {
    Body string `json:"body"`
}

type Message struct {
    ID   int    `json:"id"`
    Body string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func (msg Message) toString() string {
	return fmt.Sprintf("{%v %v %v}", msg.ID, msg.Body, msg.CreatedAt)
}


func main() {
	// Init
	createDatabaseIfNotExist("message_db")
	// Open DB Connection
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/message_db")

	if err != nil {
		log.Fatal("Failed to open connection to MySQL")
	}
	defer db.Close()

	// Handle Requests
	handleRequests()
}

func createDatabaseIfNotExist(dbName string) {

	log.Printf("Creating database '%v' if not exists.", dbName)

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/message_db")
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