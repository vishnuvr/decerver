package main

import (
	"github.com/eris-ltd/decerver"
	"github.com/eris-ltd/decerver-interfaces/glue/ipfs"
	"github.com/eris-ltd/decerver-interfaces/glue/legalmarkdown"
	"github.com/eris-ltd/decerver-interfaces/glue/monk"
	//"github.com/eris-ltd/decerver-interfaces/glue/blockchaininfo"
)

func main() {
	dc := decerver.NewDeCerver()
	mjs := monkjs.NewMonkJs()
	fm := ipfs.NewIpfs()
	lmd := legalmarkdown.NewLmdModule()
	//bci := blockchaininfo.NewBlkChainInfo()

	dc.LoadModule(lmd)
	dc.LoadModule(mjs)
	dc.LoadModule(fm)
	//dc.LoadModule(bci)

	dc.Init()

	//Run decerver
	dc.Start()
}
