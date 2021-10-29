package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getAllMsgs(c *gin.Context) {
	messages := GetAllMessages()
	if len(messages) > 0 {
		c.IndentedJSON(http.StatusOK, messages)
	} else {
		c.Data(http.StatusNotFound, "text/plain", []byte("[]"))
	}
}

func createNewUser(c *gin.Context) {
	name := c.Param("name")
	if strings.TrimSpace(name) == "" {
		c.String(http.StatusBadRequest, "User name can not be empty!")
		return
	}
	user := &User{
		Name: name,
	}
	values := c.Request.URL.Query()
	mailsParam := values["mails"]
	if len(mailsParam[0]) > 0 {
		mails := strings.Split(mailsParam[0], ",")
		user.Emails = mails
	}

	createdUser, result := NewUser(*user)
	if result {
		c.IndentedJSON(http.StatusCreated, createdUser)
	} else {
		c.String(http.StatusInternalServerError, "User created error!")
	}
}

func createMessage(c *gin.Context) {
	var message Message

	// Call BindJSON to bind the received JSON
	if err := c.BindJSON(&message); err != nil {
		c.String(http.StatusBadRequest, "Message created failed due to %s", err)
		return
	}
	if message.Time.IsZero() {
		message.Time = time.Now()
	}

	if NewMessage(message) {
		c.IndentedJSON(http.StatusCreated, message)
	} else {
		c.String(http.StatusInternalServerError, "Message created failed!")
	}

}
