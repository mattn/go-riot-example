package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/genmai"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type Todo struct {
	Id        int64     `db:"pk" column:"id" json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

func (todo *Todo) BeforeInsert() error {
	todo.CreatedAt = time.Now()
	return nil
}

func main() {
	db, err := genmai.New(&genmai.SQLite3Dialect{}, "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.CreateTableIfNotExists(&Todo{}); err != nil {
		log.Fatal(err)
	}

	goji.Get("/api/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		if err := db.Select(&todos, db.Where("id", "=", c.URLParams["id"])); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		todo := todos[0]
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	goji.Post("/api/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		if err := db.Select(&todos, db.Where("id", "=", c.URLParams["id"])); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		todo := todos[0]
		todo.Title = r.PostFormValue("title")
		todo.Done = r.PostFormValue("done") == "true"
		if _, err := db.Update(&todo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	goji.Post("/api", func(c web.C, w http.ResponseWriter, r *http.Request) {
		todo := &Todo{
			Title: r.PostFormValue("title"),
		}
		_, err := db.Insert(todo)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	goji.Get("/api", func(c web.C, w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		if err := db.Select(&todos); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&todos); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
	goji.Get("/*", http.FileServer(http.Dir("./assets")))
	goji.Serve()
}
