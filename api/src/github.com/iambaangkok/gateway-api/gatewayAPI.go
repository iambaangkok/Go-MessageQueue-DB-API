package gatewayAPI

import (
	"config"
	"context"
	"database/sql"
	"dtos"
	"fmt"
	"log"
	"models"
	"net/http"
	"utils"

	"encoding/json"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
)

var db *sql.DB
var conn *kafka.Conn

func queryRandomMessage(res http.ResponseWriter, req *http.Request, params martini.Params) []byte {

	var msg models.Message

	err := db.QueryRow(
		`SELECT *
		FROM messages
		ORDER BY RAND()
		LIMIT 1`).
		Scan(&msg.ID, &msg.Body, &msg.CreatedAt)
  
	if err != nil && err != sql.ErrNoRows {
	  log.Println(err)
	}
	log.Print(msg.ToString())

	return utils.MsgToJson(msg)
}

func queryMessageById(id int64) models.Message {

	var msg models.Message

	err := db.QueryRow(
		`SELECT *
		FROM messages m
		WHERE m.id = ?
		LIMIT 1`, id).
		Scan(&msg.ID, &msg.Body, &msg.CreatedAt)
  
	if err != nil && err != sql.ErrNoRows {
	  log.Println(err)
	}
	log.Print(msg.ToString())

	return msg
}

func queryAllMessage(res http.ResponseWriter, req *http.Request, params martini.Params) []byte {

	rows, err := db.Query(
	`SELECT *
	FROM messages
	`)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	  }
	
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var msg models.Message
        if err := rows.Scan(&msg.ID, &msg.Body, &msg.CreatedAt); err != nil {
            return utils.CompressToJsonBytes(messages)
        }
        messages = append(messages, msg)
	}

	return utils.CompressToJsonBytes(messages)
}


func addMessageKafka(res http.ResponseWriter, req *http.Request, params martini.Params) []byte {
	decoder := json.NewDecoder(req.Body)

	var dto dtos.MessageInputDTO
	err := decoder.Decode(&dto)
	if err != nil {
		panic(err)
	}

	msg := utils.CompressToJsonBytes(dto)

	_, err = conn.Write(msg)
	if err != nil {
		panic(err)
	}
	

	return utils.CompressToJsonBytes("Message added through Kafka")
}

func addMessage(res http.ResponseWriter, req *http.Request, params martini.Params) []byte {
	decoder := json.NewDecoder(req.Body)

	var dto dtos.MessageInputDTO
	err := decoder.Decode(&dto)
	if err != nil {
		panic(err)
	}

	insertStatement := fmt.Sprintf(
		`INSERT INTO messages (body)
		VALUES ("%v")`, dto.Body)

	_, err = db.ExecContext(context.Background(), insertStatement)
	if err != nil {
		log.Println(err)
	}


	return utils.CompressToJsonBytes("Message added directly to database")
}

func handleRequests() {
	m := martini.Classic()

	m.Get("/messages/random", queryRandomMessage)
	m.Get("/messages/all", queryAllMessage)
	m.Post("/messages/add", addMessageKafka)
	m.Post("/messages/add-direct", addMessage)
	m.Run()
}

func main() {
	/// Init
	// Init Kafka
	cfg := config.KafkaConnCfg {
		Url:   "localhost:9092",
		Topic: "message.topic",
	}
	conn = utils.KafkaConn(cfg)
	utils.KafkaCreateTopicIfNotExist(conn, cfg)
	defer utils.KafkaCloseConn(conn)

	

	// Init DB
	utils.CreateDatabaseIfNotExist("message_db")
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