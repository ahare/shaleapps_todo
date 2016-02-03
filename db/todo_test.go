package db

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	os.Remove(dbPath)
	Load()
	defer Close()
	os.Exit(m.Run())
}

func TestNewTodo(t *testing.T) {
	text := randText()
	todo := NewTodo(text)
	if todo.Text != text {
		t.Fatalf("text was %q, expected %q", todo.Text, text)
	}
	if todo.Done {
		t.Fatal("done was true, expected false")
	}
}

func TestCreateTodo(t *testing.T) {
	text := randText()
	todo := NewTodo(text)
	err := todo.Save()
	if err != nil {
		t.Fatalf("save failed: %q", err)
	}
	if todo.Id == 0 {
		t.Fatal("id wasn't updated")
	}
	foundTodo, err := FindTodoById(todo.Id)
	if err != nil {
		t.Fatalf("find by id failed: %q", err)
	}
	if !reflect.DeepEqual(todo, foundTodo) {
		t.Fatal("todo was not inserted correctly")
	}
}

func TestUpdateTodo(t *testing.T) {
	todo := NewTodo("test")
	todo.Save()

	text := randText()
	todo.Text = text
	todo.Done = true
	todo.Save()

	foundTodo, err := FindTodoById(todo.Id)
	if err != nil {
		t.Fatalf("find by id failed: %q", err)
	}
	if !reflect.DeepEqual(todo, foundTodo) {
		t.Fatal("todo was not updated correctly")
	}
}

func TestDeleteTodo(t *testing.T) {
	todo := NewTodo(randText())
	todo.Save()
	foundTodo, err := FindTodoById(todo.Id)
	if err != nil {
		t.Fatalf("find by id failed: %q", err)
	}
	if !reflect.DeepEqual(todo, foundTodo) {
		t.Fatal("todo was not inserted correctly")
	}

	err = DeleteTodo(todo.Id)
	if err != nil {
		t.Fatalf("delete failed: %q", err)
	}

	deletedTodo, err := FindTodoById(todo.Id)
	if err != nil {
		t.Fatal(err)
	}
	if deletedTodo != nil {
		t.Fatal("delete was not successful - todo still found by ID")
	}
}

func TestFindTodosByText(t *testing.T) {
	for _, name := range []string{"dog", "dogged", "doggy", "cat"} {
		todo := NewTodo(name)
		todo.Save()
	}

	todos, err := FindTodosByText("dog")
	if err != nil {
		t.Fatal(err)
	}
	if len(todos) != 3 {
		t.Fatal("wrong number of todos returned")
	}
}

func TestFindTodosByDone(t *testing.T) {
	db.Exec("DELETE FROM todos;")

	for i, name := range []string{"1", "2", "3", "4", "5"} {
		todo := NewTodo(name)
		todo.Done = i%2 == 0
		todo.Save()
	}

	dones, err := FindTodosByDone(true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dones)
	if len(dones) != 3 {
		t.Fatal("wrong number of done todos returned")
	}

	notDones, err := FindTodosByDone(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(notDones) != 2 {
		t.Fatal("wrong number of not done todos returned")
	}
}
