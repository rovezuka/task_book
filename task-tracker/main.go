package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Добавление задачи в базу данных
func addTask(db *sql.DB, description string) error {
	_, err := db.Exec(`INSERT INTO tasks (description, done) VALUES (?, ?)`, description, false)
	return err
}

// Получение задач из базы данных
func getTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query(`SELECT * FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Description, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Обновление статуса задачи в базе данных
func updateTaskStatus(db *sql.DB, taskID int, done bool) error {
	_, err := db.Exec(`UPDATE tasks SET done = ? WHERE id = ?`, done, taskID)
	return err
}

// Удаление задачи из базы данных
func deleteTask(db *sql.DB, taskID int) error {
	_, err := db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	return err
}

func main() {

	// Инициализация базы данных
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы задач, если она не существует
	_, err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT,
		done BOOLEAN
	    )
	`)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	// Обработчик для статических файлов
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		tasks, err := getTasks(db)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tasks": tasks,
		})
	})

	router.POST("/delete/:index", func(c *gin.Context) {
		index := c.Param("index")

		i, err := strconv.Atoi(index)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = deleteTask(db, i)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/add", func(c *gin.Context) {
		newTask := c.PostForm("newTask")
		if newTask != "" {
			err := addTask(db, newTask)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/update/:index", func(c *gin.Context) {
		index := c.Param("index")
		i, err := strconv.Atoi(index)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = updateTaskStatus(db, i, true)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8080")
}
