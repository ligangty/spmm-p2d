package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleMsg(c *gin.Context) {
	path := os.Getenv("TEMPLATE_PATH")
	if path == "" {
		path = "/var/www/template.html"
	}

	msg_content := ""
	status := http.StatusOK
	contentType := "text/plain"
	template_content, err := os.ReadFile(path)
	if err != nil {
		msg_content = err.Error()
		status = http.StatusInternalServerError
	} else {
		message := c.Param("message")
		template_content_str := strings.TrimLeft(string(template_content), " ")
		if strings.HasPrefix(template_content_str, "<html>") {
			contentType = "text/html"
		} else if strings.HasPrefix(template_content_str, "{") {
			contentType = "application/json"
		}
		msg_content = strings.ReplaceAll(template_content_str, "${message}", message)
	}

	c.Data(status, contentType, []byte(msg_content))
}
