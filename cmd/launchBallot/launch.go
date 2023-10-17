package main

import (
	"fmt"

	agt "github.com/adrsimon/voting-system-ia04/agt"
)

func main() {
	server := agt.NewServerRest(":8080")
	server.Start()
	fmt.Scanln()
}
