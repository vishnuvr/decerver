package main 

import (
	"github.com/eris-ltd/deCerver"
	"github.com/eris-ltd/deCerver/testmodule"
)

func main() {
	dc := deCerver.NewDeCerver()
	dc.Init()
	testMod := &testmodule.TestModule{}
	dc.AddModule(testMod.Name(),testMod)
	// Run deCerver
	dc.Run()
}

