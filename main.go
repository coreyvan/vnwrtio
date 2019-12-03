package main

import (
	"github.com/coreyvan/vnwrtio/server"
)

func main() {
	s := server.Server{}

	s.HandleSignals()
}
