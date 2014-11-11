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
	dc.LoadModule(mm)
	dc.LoadModule(fm)
	
	dc.Init()
	// Run deCerver
	dc.Start()
}
