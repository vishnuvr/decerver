package main

import (
	"github.com/eris-ltd/deCerver/administrator/ethmanager"
	//"github.com/eris-ltd/deCerver/administrator/profiler"
	"github.com/eris-ltd/deCerver/administrator/ponos"
	"github.com/eris-ltd/deCerver/administrator/server"
	"github.com/eris-ltd/thelonious/monk"
)

func main() {
	
	// Don't need anything but the ethereum here
	ec := monk.NewEth(nil)
	ec.Init()
	srpcf := ethmanager.NewEthereumSRPCFactory(ec.Ethereum)
	
	prpcf := ponos.NewPonosSRPCFactory(ec.Ethereum)
	
	ws := server.NewWebServer(10)
	// Pass the various rpc services.
	ws.RegisterSocketRPCServices(srpcf, prpcf)
	
	ws.Start()
}
