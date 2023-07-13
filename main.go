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
	Status  string
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

	r.PUT("/todo/:id/complete/:status", func(ctx *gin.Context) {
		status := ctx.Param("status")
		log.Printf("STATUS\n\n%s\n\n", status)
		todo, hasErr := getByID(db, ctx.Param("id"))
		if hasErr {
			ctx.HTML(http.StatusNotFound, "todo.html", gin.H{
				"todo": todo,
			})
			return
		}

		// toggle status
		if status == "checked" {
			status = "pending"
		} else {
			status = "checked"
		}
		todo.Status = status

		updateTODO(db, todo)
		log.Println(todo)

		ctx.HTML(http.StatusOK, "todo.html", gin.H{
			"todo": todo,
		})
	})

	r.Run(":8080")
}

func updateTODO(db *sql.DB, todo TODO) {
	var cmp int
	if todo.Status == "checked" {
		cmp = 1
	}
	_, err := db.Exec(`UPDATE todos
	SET title = ?,
	content = ?,
	completed = ?
	WHERE rowid = ?
	`, todo.Title, todo.Content, cmp, todo.ID)

	if err != nil {
		log.Printf("error updating todo: %s", err)
	}
}

func getByID(db *sql.DB, id string) (TODO, bool) {
	log.Println("ID TO PRINT", id)
	var (
		idVal     int
		title     string
		content   string
		completed int
		cmpStr    string
	)
	err := db.QueryRow("SELECT rowid, title, content, completed FROM todos WHERE rowid = ?", id).
		Scan(&idVal, &title, &content, &completed)
	if err != nil {
		log.Printf("error getting todo: %s", err)
		return TODO{}, true
	}

	if completed > 0 {
		cmpStr = "checked"
	} else {
		cmpStr = "pending"
	}

	return TODO{
		ID:      idVal,
		Title:   title,
		Content: content,
		Status:  cmpStr,
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

	rows, err := db.Query("SELECT rowid, title, content, completed FROM todos")
	if err != nil {
		return todos
	}
	defer rows.Close()

	for rows.Next() {
		var (
			todo TODO
			cmp  int
		)
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &cmp); err != nil {
			continue
		}
		if cmp != 0 {
			log.Printf("checked todo: %s, %s", todo.Title, todo.Content)
			todo.Status = "checked"
		} else {
			todo.Status = "pending"
		}
		todos = append(todos, todo)
	}

	return todos
}

func createDatabase(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS todos (title TEXT NOT NULL, content TEXT NOT NULL, completed INT DEFAULT 0);")
	if err != nil {
		log.Fatalf("error creating table: %s", err)
	}
	_, err = db.Exec(`INSERT INTO todos (title, content, completed) VALUES
	("title one", "content 1", 0),
	("title two", "content 2", 1),
	("title three", "content three", 0);`)
	if err != nil {
		log.Fatalf("error inserting data: %s", err)
	}
}
