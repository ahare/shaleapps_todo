package db

import (
	"database/sql"
	"fmt"
)

type Todo struct {
	Id   int64  `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func init() {
	queries["todos.insert"] = `INSERT INTO todos(text, done) VALUES(?, ?);`
	queries["todos.update"] = `UPDATE todos SET text = ?, done = ? WHERE id = ?;`
	queries["todos.findById"] = `SELECT text, done FROM todos WHERE id = ?;`
	queries["todos.findByText"] = `SELECT id, text, done FROM todos WHERE text LIKE ?;`
	queries["todos.findByDone"] = `SELECT id, text, done FROM todos WHERE done = ?;`
	queries["todos.all"] = `SELECT id, text, done FROM todos;`
	queries["todos.delete"] = `DELETE FROM todos WHERE id = ?;`
}

func NewTodo(text string) *Todo {
	return &Todo{Text: text, Done: false}
}

func FindTodoById(id int64) (*Todo, error) {
	find, err := getPrep("todos.findById")
	if err != nil {
		return nil, err
	}

	var (
		text sql.NullString
		done sql.NullBool
	)

	err = find.QueryRow(id).Scan(&text, &done)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &Todo{id, text.String, done.Bool}, nil
	}
}

func FindTodosByText(text string) ([]*Todo, error) {
	return findTodos("todos.findByText", fmt.Sprintf("%%%s%%", text))
}

func FindTodosByDone(done bool) ([]*Todo, error) {
	return findTodos("todos.findByDone", done)
}

func AllTodos() ([]*Todo, error) {
	return findTodos("todos.all", []interface{}{}...)
}

func findTodos(queryName string, values ...interface{}) ([]*Todo, error) {
	find, err := getPrep(queryName)
	if err != nil {
		return nil, err
	}

	rows, err := find.Query(values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []*Todo{}
	for rows.Next() {
		var (
			id   sql.NullInt64
			text sql.NullString
			done sql.NullBool
		)
		rows.Scan(&id, &text, &done)
		todos = append(todos, &Todo{id.Int64, text.String, done.Bool})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func DeleteTodo(id int64) error {
	delete, err := getPrep("todos.delete")
	if err != nil {
		return err
	}
	_, err = delete.Exec(id)
	return err
}

func (t *Todo) Save() error {
	if t.Id == 0 {
		return t.insert()
	}
	return t.update()
}

func (t *Todo) insert() error {
	insert, err := getPrep("todos.insert")
	if err != nil {
		return err
	}

	res, err := insert.Exec(t.Text, t.Done)
	if err != nil {
		return err
	}

	t.Id, err = res.LastInsertId()
	return err
}

func (t *Todo) update() error {
	update, err := getPrep("todos.update")
	if err != nil {
		return err
	}

	_, err = update.Exec(t.Text, t.Done, t.Id)
	if err != nil {
		return err
	}

	return nil
}
