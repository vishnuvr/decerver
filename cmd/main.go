package main

import (
	"github.com/eris-ltd/deCerver"
	"github.com/eris-ltd/thelonious/monk"
	"github.com/eris-ltd/glululemon/ipfs"
)

func main() {
	
	dc := deCerver.NewDeCerver()
	mm := monk.NewMonk(nil)
	fm := ipfs.NewIpfs()
	dc.AddModule(mm)
	dc.AddModule(fm)
	
	dc.Init()
	// Run deCerver
	dc.Start()
}
