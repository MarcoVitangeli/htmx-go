package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type TODO struct {
	Title   string
	Content string
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())

	r.LoadHTMLGlob("./html/*.html")
	r.StaticFile("css/pico.min.css", "./pico-1.5.9/css/pico.min.css")

	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		log.Fatalf("error creating database: %s", err)
	}
	defer db.Close()

	createDatabase(db)

	r.GET("/", func(ctx *gin.Context) {
		todos := getTodos(db)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"timestamp": time.Now().UTC().String(),
			"todos":     todos,
		})
	})

	r.Run(":8080")
}

func getTodos(db *sql.DB) []TODO {
	var todos []TODO

	rows, err := db.Query("SELECT title, content FROM todos")
	if err != nil {
		return todos
	}
	defer rows.Close()

	for rows.Next() {
		var todo TODO
		if err := rows.Scan(&todo.Title, &todo.Content); err != nil {
			continue
		}
		todos = append(todos, todo)
	}

	return todos
}

func createDatabase(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS todos (title TEXT NOT NULL, content TEXT NOT NULL);")
	if err != nil {
		log.Fatalf("error creating table: %s", err)
	}
	_, err = db.Exec(`INSERT INTO todos (title, content) VALUES
	("title one", "content 1"),
	("title two", "content 2"),
	("title three", "content three");`)
	if err != nil {
		log.Fatalf("error inserting data: %s", err)
	}
}
