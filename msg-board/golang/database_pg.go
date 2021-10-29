package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func createSchema() error {
	db := getDB()
	defer db.Close()
	models := []interface{}{
		(*User)(nil),
		(*Message)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewUser(user User) (User, bool) {
	db := getDB()
	defer db.Close()
	_, err := db.Model(&user).Insert()
	if err != nil {
		fmt.Printf("Error to create new user: %s, error: %s", user.Name, err)
		return *new(User), false
	}

	return user, true
}

func GetUserById(db *pg.DB, userId int64) User {
	var user = new(User)
	sql := `SELECT * FROM users WHERE id = ?`
	_, err := db.QueryOne(user, sql, userId)
	panicErr(err)
	return *user
}

func GetAllUsers(db *pg.DB) []User {
	var users = new([]User)
	sql := `SELECT * FROM users`
	_, err := db.Query(users, sql)
	if err != nil {
		panic(err)
	}
	return *users
}

func GetMessageById(db *pg.DB, msgId int64) Message {
	var msg = new(Message)
	sql := `SELECT * FROM messages WHERE id = ?`
	_, err := db.QueryOne(msg, sql, msgId)
	panicErr(err)
	return *msg
}

func GetAllMessages() []Message {
	db := getDB()
	defer db.Close()
	var msgs = new([]Message)
	sql := `SELECT * FROM messages`
	_, err := db.Query(msgs, sql)
	if err != nil {
		panic(err)
	}
	return *msgs
}

func GetMessagesByUser(user User) []Message {
	db := getDB()
	defer db.Close()
	var msgs = new([]Message)
	sql := `SELECT * FROM messages WHERE user_id = ?`
	_, err := db.Query(msgs, sql, user.Id)
	panicErr(err)
	return *msgs
}

func NewMessage(message Message) bool {
	db := getDB()
	defer db.Close()

	_, err := db.Model(&message).Insert()
	if err != nil {
		fmt.Printf("Failed to create message due to %s", err)
		return false
	}
	return true
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getDB() *pg.DB {
	return pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Addr:     "localhost:5432",
		Database: "msgboard",
	})
}
