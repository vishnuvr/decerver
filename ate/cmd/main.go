package main

import (
	"fmt"
	"github.com/obscuren/otto"
)

type TestStruct struct {
	Val0 string
	Val1 int
}

func (this *TestStruct) GetVal0() string {
	return this.Val0
}

func main() {
	vm := otto.New()
	theVal := &TestStruct{}
	theVal.Val0 = "testString"
	theVal.Val1 = 1
	// var stuff interface{}
	// stuff = theVal
	vm.Set("test", theVal)
	val, err := vm.Run("test.GetVal0()")

	if err != nil {
		fmt.Println(err.Error())
	}
	
	fmt.Printf("Value: %v\n", val)
}
