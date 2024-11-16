package main

import (
	"transactions/server"

	_ "github.com/lib/pq"
)

func main() {
	server := server.NewServer()
	server.Start()
}
