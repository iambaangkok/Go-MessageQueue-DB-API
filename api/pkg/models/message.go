package models

import (
	"fmt"
)

type Message struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func (msg Message) ToString() string {
	return fmt.Sprintf("{%v %v %v}", msg.ID, msg.Body, msg.CreatedAt)
}