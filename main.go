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
	ID      int
	Title   string
	Content string
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())

	r.LoadHTMLGlob("./html/*.html")

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

	r.POST("/todo", func(ctx *gin.Context) {
		title := ctx.Request.FormValue("title")
		desc := ctx.Request.FormValue("description")

		createTODO(db, title, desc)

		todos := getTodos(db)
		ctx.HTML(http.StatusOK, "todos.html", gin.H{
			"todos": todos,
		})
	})

	r.DELETE("/todo/:id", func(ctx *gin.Context) {
		_, err := db.Exec("DELETE FROM todos WHERE rowid = ?", ctx.Param("id"))
		if err != nil {
			log.Printf("error deleting todo: %s", err)
			ctx.Status(http.StatusInternalServerError)
		}
	})

	r.Run(":8080")
}

func getByID(db *sql.DB, id string) (TODO, bool) {
	log.Println("ID TO PRINT", id)
	var (
		idVal     int
		title     string
		content   string
		completed int
	)
	err := db.QueryRow("SELECT rowid, title, content, completed FROM todos WHERE rowid = ?", id).
		Scan(&idVal, &title, &content, &completed)
	if err != nil {
		log.Printf("error getting todo: %s", err)
		return TODO{}, true
	}

	return TODO{
		ID:      idVal,
		Title:   title,
		Content: content,
	}, false
}

func createTODO(db *sql.DB, title, desc string) {
	_, err := db.Exec("INSERT INTO todos (title, content) VALUES (?,?)", title, desc)
	if err != nil {
		log.Printf("error inserting TODO: %s", err)
	} else {
		log.Println("success when inserting TODO")
	}
}

func getTodos(db *sql.DB) []TODO {
	var todos []TODO

	rows, err := db.Query("SELECT rowid, title, content FROM todos")
	if err != nil {
		return todos
	}
	defer rows.Close()

	for rows.Next() {
		var (
			todo TODO
		)
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Content); err != nil {
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
