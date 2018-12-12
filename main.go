package main

import (
	"github.com/Popcore/verisart/pkg/server"
)

func main() {
	s := server.New(":9091")
	s.Start()
}
