package main

import (
	_ "github.com/lib/pq"
	"github.com/shashimalcse/tiny-is/internal/server"
)

func main() {

	server.StartServer()
}
