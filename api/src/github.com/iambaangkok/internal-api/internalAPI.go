package internalAPI

import (
	"config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"utils"

	"github.com/segmentio/kafka-go"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func insertMessage(kMsg kafka.Message) []byte {

	dto := utils.KafkaMessageToMessageInputDTO(kMsg) 

	insertStatement := fmt.Sprintf(
		`INSERT INTO messages (body)
		VALUES ("%v")`, dto.Body)

	_, err := db.ExecContext(context.Background(), insertStatement)
	if err != nil {
		log.Println(err)
	}

	return utils.CompressToJsonBytes("Success")
}


func consumeMessages(conn *kafka.Conn) {
	for {
		message, err := conn.ReadMessage(10e3)
		if err != nil {
			break
		}

		fmt.Println(string(message.Value))
		insertMessage(message)
	}
}

func main() {
	/// Init
	// Init Kafka
	cfg := config.KafkaConnCfg {
		Url:   "localhost:9092",
		Topic: "message.topic",
	}
	conn := utils.KafkaConn(cfg)
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

	// Start Consuming
	consumeMessages(conn)
}
