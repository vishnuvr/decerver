package main

import (
	"github.com/eris-ltd/decerver"
	"github.com/eris-ltd/decerver-modules/ipfs"
	"github.com/eris-ltd/decerver-modules/legalmarkdown"
	"github.com/eris-ltd/decerver-modules/monk"
	//"github.com/eris-ltd/modules/blockchaininfo"
)

func main() {
	dc := decerver.NewDeCerver()
	mjs := monkjs.NewMonkJs()
	fm := ipfs.NewIpfsModule()
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
