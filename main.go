package main

import (
	"github.com/popcore/verisart_exercise/pkg/server"
)

func main() {
	s := server.New(":9091")
	s.Start()
}
