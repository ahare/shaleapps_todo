package db

import (
	"strconv"
	"testing"
	"time"
)

func init() {
	dbPath = "../todo_test.db"
}

func randText() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func TestGetPrep(t *testing.T) {
	name := randText()
	queries[name] = `SELECT * FROM todos;`
	_, err := getPrep(name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPrepNoQuery(t *testing.T) {
	_, err := getPrep(randText())
	if err == nil {
		t.Fatal("should have returned error for missing query")
	}
}
