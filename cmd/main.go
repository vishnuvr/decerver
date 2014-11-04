package main

import (
	"github.com/eris-ltd/deCerver"
	"github.com/eris-ltd/deCerver-interfaces/monk"
)

func main() {
	dc := deCerver.NewDeCerver()
	dc.Init()
	mm := monk.NewMonkModule()
	dc.AddModule(mm)
	// Run deCerver
	dc.Start()
}
