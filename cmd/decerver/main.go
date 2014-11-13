package main

import (
	"github.com/eris-ltd/decerver"
	"github.com/eris-ltd/thelonious/monk"
	"github.com/eris-ltd/glululemon/ipfs"
)

func main() {
	
	dc := decerver.NewDeCerver()
	mm := monk.NewMonk(nil)
	fm := ipfs.NewIpfs()
	dc.LoadModule(mm)
	dc.LoadModule(fm)
	
	dc.Init()
	// Run decerver
	dc.Start()
}
