package main

import (
	"fmt"
	"time"
)

type Message struct {
	Id      int64     `json:"id" pg:"id,pk"`
	UserId  int64     `json:"uid" pg:"user_id"`
	Content string    `json:"content" pg:"content"`
	Time    time.Time `json:"time" pg:"time"`
}

func (s Message) String() string {
	return fmt.Sprintf("Message<%d %s %d>", s.Id, s.Content, s.UserId)
}
