package main

import (
	"github.com/eris-ltd/decerver"
	"github.com/eris-ltd/decerver-interfaces/glue/ipfs"
	"github.com/eris-ltd/decerver-interfaces/glue/monk"
)

func main() {
	dc := decerver.NewDeCerver()
	mjs := monkjs.NewMonkJs()
	fm := ipfs.NewIpfs()

	dc.LoadModule(mjs)
	dc.LoadModule(fm)

	dc.Init()
	
	//Run decerver
	dc.Start()

}
