package main

import (
	"github.com/joaosoft/monitor"
)

func main() {
	m, err := monitor.NewMonitor()
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}
}
