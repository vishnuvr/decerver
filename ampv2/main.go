package main

import (
	"fmt"
	"github.com/eris-ltd/deCerver/scriptengine"
	"github.com/eris-ltd/thelonious/monk"
	"os"
)

func main() {
	ec := monk.NewEth(nil)
	ec.Init()
	se := scriptengine.NewScriptEngine(ec)
	
	se.AddModel("GenDoug","genDoug.js")
	
	res, errRun := se.RunAction("GenDoug","getContracts", nil)
	if errRun != nil {
		fmt.Println(errRun.Error())
		os.Exit(0)
	}
		
	fmt.Printf("Result: %s\n",res)
}

/*
func main() {

	// New vm
	vm := otto.New()

	// Simple example of how to run script (console used stdout)
	vm.Run(`
	    abc = 2 + 2;
	    console.log("The value of abc is " + abc);
	`)

	// Define GetStorageAt(account,addr), so that it is accessible from within the vm.
	// Notice the vm is already started and running. This is not part of vm initialization
	// or something, but can be done at any time.
	vm.Set("GetStorageAt", func(call otto.FunctionCall) otto.Value {
		account, _ := call.Argument(0).ToString()
		address, _ := call.Argument(1).ToString()
		// Bind it to the Go "FakeGetStorageAt" function, which would normally call eth.
		result, _ := vm.ToValue(FakeGetStorageAt(account, address))
		return result
	})

	// Add some stuff. We could have loaded a script file instead based on some hash found
	// in a contract.
	vm.Run(`
		// Not doing this all out with closures etc. Just for demonstrating.
	    gendoug = {
	    	"addr" : "0x00000",
	    	"db" : "0x30",
	    };
	    
	    gendoug.getDeadbeef = function (){
	        return GetStorageAt(gendoug.addr,gendoug.db)
	    };
	    
	`)

	// Make a "blockchain" query based on an incoming api call.
	result, _ := vm.Run(`
    	result = gendoug.getDeadbeef()
    `)
	fmt.Printf("Result: %v\n", result)
	
	// Now extend gendoug. This should of course not be done in this way.
	vm.Run(`
	    gendoug.scriptHashAddr = "0x19";
	    gendoug.getScriptHash = function (){
	        return GetStorageAt(gendoug.addr,gendoug.scriptHashAddr);
	    };
	`)
	
	// Now have at it
	result, _ = vm.Run(`
    	result = gendoug.getScriptHash();
    `)
	
	fmt.Printf("Result: %v\n", result)
}

func FakeGetStorageAt(account, address string) string {
	if address == "0x30" {
		return "0xdeadbeef"
	} else if address == "0x19" {
		return "0x47569283769254675942837569"
	}
	return "0x"
}
*/
