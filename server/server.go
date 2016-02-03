package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shaleapps_todo/db"
	"strconv"
	"strings"
)

func Start() {
	r := &router{}
	r.HandleFunc("/todos/\\d+", "DELETE", deleteTodoHandler)
	r.HandleFunc("/todos/.*", "PUT", updateTodoHandler)
	r.HandleFunc("/todos", "POST", createTodoHandler)
	r.HandleFunc("/todos/\\d+", "GET", findTodoByIdHandler)
	r.HandleFunc("/todos?.*", "GET", findTodosHandler)

	fmt.Println("ToDo server listening on :8080...")
	http.ListenAndServe(":8080", r)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	idStr := urlParts[len(urlParts)-1]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = db.DeleteTodo(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	saveTodoHandler(http.StatusOK, w, r)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	saveTodoHandler(http.StatusCreated, w, r)
}

func saveTodoHandler(successStatus int, w http.ResponseWriter, r *http.Request) {
	t, err := saveTodo(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(successStatus)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func saveTodo(r *http.Request) (*db.Todo, error) {
	decoder := json.NewDecoder(r.Body)
	var t db.Todo
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	err = t.Save()
	return &t, err
}

func findTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	idStr := urlParts[len(urlParts)-1]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := db.FindTodoById(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func findTodosHandler(w http.ResponseWriter, r *http.Request) {
	var (
		todos []*db.Todo
		err   error
	)
	values := r.URL.Query()

	doneStr := values.Get("done")
	if len(doneStr) > 0 {
		done, err := strconv.ParseBool(doneStr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		todos, err = db.FindTodosByDone(done)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	text := values.Get("text")
	if len(text) > 0 {
		todos, err = db.FindTodosByText(text)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
