package main

import (
	"fmt"

	server "github.com/adrsimon/voting-system-ia04/server"
)

func main() {
	server := server.NewVoteServer(":8080")
	server.Start()
	fmt.Scanln()
}
