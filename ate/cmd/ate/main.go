package main

import (
	"fmt"
	"github.com/obscuren/otto"
	//"github.com/obscuren/otto/parser"
)

type TestStruct struct {
	Val0 string       `json:"val0"`
	Val1 TestStruct2 `json:"val1"`
	Val2 []string
}

func (this *TestStruct) GetVal0() string {
	return this.Val0
}

func (this *TestStruct) GetVal1() TestStruct2 {
	return this.Val1
}

func (this *TestStruct) GetVal2() []string {
	return this.Val2
}

type TestStruct2 struct {
	Val string `json:"val"`
}

func main() {
	ts := &TestStruct{"test",
		TestStruct2{"testInner"}, 
		[]string{"testInString"},
	}
	/*
	valMap := make(map[string]interface{})
	valMap["testString"] = "test"
	valMap["testStringArray"] = []string{"testing"}
	
	intMap := make(map[string]interface{})
	intMap["test"] = 55
	
	valMap["testIntMap"] = intMap
	*/
	vm := otto.New()

	vm.Set("testStruct", ts)
	val, err := vm.Run("testStruct.GetVal1()")

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%v\n", val)
	}
}
