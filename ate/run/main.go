package main

import (
	//"fmt"
	"github.com/eris-ltd/deCerver/ate"
	"github.com/eris-ltd/thelonious/monk"
)

func main() {
	ec := monk.NewEth(nil)
	ec.Init()
	
	//se := scriptengine.NewScriptEngine(ec)
	
	//se.AddModel("GenDoug","genDoug.js")
	
	//res, errRun := se.RunAction("GenDoug","getContracts", nil)
	//if errRun != nil {
	//	fmt.Println(errRun.Error())
	//	os.Exit(0)
	//}
	
	//fmt.Printf("Result: %s\n",res)
	tp := ate.NewTreeParser(ec)
	//tp.ParseTree("ponos","0x99883cb40909f59f8082128afd4b67a2ca2b206c")
	tree, _ := tp.ParseTree("ponos","0x99883cb40909f59f8082128afd4b67a2ca2b206c")
	tp.FlattenTree(tree)
	
}