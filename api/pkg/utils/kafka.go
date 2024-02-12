package utils

import (
	"context"
	"encoding/json"
	"log"
	"main/api/config"
	"main/api/dtos"
	"time"

	// "github.com/Rayato159/go-simple-kafka/config"

	"github.com/segmentio/kafka-go"
)

func KafkaConn(cfg config.KafkaConnCfg) *kafka.Conn {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Url, cfg.Topic, 0)
 	if err != nil {
  		panic(err.Error())
 	}
 	return conn
}

func KafkaCloseConn(conn *kafka.Conn) {
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func KafkaCreateTopicIfNotExist(conn *kafka.Conn, cfg config.KafkaConnCfg) {
	   
	if !IsTopicAlreadyExists(conn, cfg.Topic) {
		topicConfigs := []kafka.TopicConfig{
		 	{
				Topic:             cfg.Topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
		 	},
		}
	  
		err := conn.CreateTopics(topicConfigs...)
		if err != nil {
			panic(err.Error())
		}
	}
		
}

func MessageInputDTOToKafkaMessage(dto dtos.MessageInputDTO) kafka.Message {
	return kafka.Message{
		Value: CompressToJsonBytes(&dto),
	}
}

func KafkaMessageToMessageInputDTO(msg kafka.Message) dtos.MessageInputDTO {
	var dto dtos.MessageInputDTO

	err := json.Unmarshal(msg.Value, &dto)
    if err != nil {
		log.Println(err)
    }

	return dto
}

func Write(conn *kafka.Conn, msg kafka.Message) {
	// Set timeout
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.Write(msg.Value)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func IsTopicAlreadyExists(conn *kafka.Conn, topic string) bool {
 	partitions, err := conn.ReadPartitions()
 	if err != nil {
  		panic(err.Error())
 	}

 	for _, p := range partitions {
  		if p.Topic == topic {
   			return true
  		}
 	}
 	return false
}