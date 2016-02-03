package main

import (
	"shaleapps_todo/db"
	"shaleapps_todo/server"
)

func main() {
	db.Load()
	defer db.Close()

	server.Start()
}
