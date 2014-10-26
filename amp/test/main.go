package main

import (
    "io/ioutil"
	"fmt"
	"github.com/eris-ltd/deCerver/amp"
	"github.com/eris-ltd/thelonious/monk"
	"github.com/eris-ltd/thelonious/ethchain"
	"os"
	"encoding/hex"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test action (example: test deadbeef)")
		os.Exit(0)
	}
	createModels()
	ec := monk.NewEth(nil)
	ec.Init()
	parser := amp.NewActionModelParser(ec)
	
	contract := "0x" + hex.EncodeToString(ethchain.GENDOUG)
	fmt.Printf("ACTION MODEL TEST: (contract: %s)\n",contract)
	err := parser.Initialize(contract, "comments")
	
	if err != nil {
		panic(err)
	}
	
	res := parser.PerformAction(os.Args[1], []string{"0"})
	fmt.Printf("Result (success, return): %v\n",res)
}

func createModels() {
	bytes, err := ioutil.ReadFile("gendoug.json")
	if err != nil {
		panic(err)
	}
	amp.ActionModels = make(map[string][]byte)
	amp.ActionModels["comments"] = bytes
	
}
