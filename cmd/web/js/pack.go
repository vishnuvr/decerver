package main 

import (
	"io/ioutil"
	"fmt"
)

var bts []byte = make([]byte,0)

func main() {
	
	addFile("./jquery-1.11.0.js")
	addFile("./plugins/dataTables/jquery.dataTables.js")
	addFile("./bootstrap.min.js")
	addFile("./plugins/dataTables/dataTables.bootstrap.js")
	addFile("./plugins/enscroll/enscroll-0.6.1.min.js")
	addFile("./plugins/terminal/jquery.terminal-0.8.8.js")
	addFile("./plugins/terminal/jquery.mousewheel-min.js")
	addFile("./app/rpc.js")
	addFile("./app/BigInteger.js")
	addFile("./app/ethString.js")
	addFile("./app/ethutil.js")
	addFile("./app/ethrpc.js")
	addFile("./app/ethadmin.js")
	addFile("./app/ea_eventhandlers.js")
	fmt.Println(len(bts))
	
	ioutil.WriteFile("mega.js",bts,0600)
}

func addFile(fileName string) {
	bb, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	bts = append(bts,bb...)
}