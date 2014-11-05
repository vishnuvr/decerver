package main

import (
	//"fmt"
	//"github.com/obscuren/otto"
	//"github.com/obscuren/otto/parser"
)

type TestStruct struct {
	val0 string
	val1 int
}

func (this *TestStruct) GetVal0() string {
	return this.val0
}

func (this *TestStruct) GetVal1() int {
	return this.val1
}

type Test0 interface {
	GetVal0() string
}

type Test1 interface {
	GetVal1() int
}

func main() {
	/*
	vm := otto.New()
	
	   	filename := "" // A filename is optional
	   	src := `
	       // Sample xyzzy example
	       (function(){
	           if (3.14159 > 0) {
	               console.log("Hello, World.");
	               return;
	           }

	           var xyzzy = NaN;
	           console.log("Nothing happens.");
	           return xyzzy;
	       })();
	   `



	   	// Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
	   	program, err := parser.ParseFile(nil, filename, src, 0)

	   	if err != nil {
	   		fmt.Println(err.Error())
	   	}

	   	stmts := program.Body
	   	for _ , s := range stmts {
	   		fmt.Printf("Stuff: %v\n",s)
	   	}
	
	
	testStruct := &TestStruct{}
	testStruct.val0 = "testVal"
	testStruct.val1 = 4
	
	vm.Set("TestObj", obj)
	val, err := vm.Run("TestObj.GetVal1()")
	
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Value: %v\n", val)
	*/
}
