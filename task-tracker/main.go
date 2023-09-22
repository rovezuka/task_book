package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Инициализация базы данных с использованием GORM и драйвера SQLite
	db, err := gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Миграция таблицы задач в базу данных
	db.AutoMigrate(&Task{})

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	// Обработчик для статических файлов
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		var tasks []Task
		// Обработка ошибки при получении задач
		if err := db.Find(&tasks).Error; err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tasks": tasks,
		})
	})

	router.POST("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		var task Task
		if err := db.First(&task, id).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Обработка ошибки при удалении задачи
		if err := db.Delete(&task).Error; err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/add", func(c *gin.Context) {
		newTask := c.PostForm("newTask")
		if newTask != "" {
			task := Task{
				Description: newTask,
				Done:        false, // Начальное состояние задачи
			}

			// Обработка ошибки при создании задачи
			if err := db.Create(&task).Error; err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		var task Task
		if err := db.First(&task, id).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Обновляем статус задачи
		task.Done = true

		// Обработка ошибки при обновлении задачи
		if err := db.Save(&task).Error; err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8080")
}
