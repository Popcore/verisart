package main

import (
	"github.com/popcore/verisart/pkg/server"
)

func main() {
	s := server.New(":9091")
	s.Start()
}
