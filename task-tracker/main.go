package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var tasks = []string{}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	// Обработчик для статических файлов
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"tasks": tasks,
		})
	})

	router.POST("/delete/:index", func(c *gin.Context) {
		index := c.Param("index")

		i, err := strconv.Atoi(index)
		if err == nil && i >= 0 && i < len(tasks) {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/add", func(c *gin.Context) {
		newTask := c.PostForm("newTask")
		tasks = append(tasks, newTask)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8080")
}
