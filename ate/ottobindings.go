package ate

import (
	"encoding/hex"
	"fmt"
	"github.com/obscuren/sha3"
	"github.com/robertkrimen/otto"
	//"github.com/eris-ltd/deCerver-interfaces/events"
	"math/big"
)

// Enables math operations on strings using big.Int.
var BZERO *big.Int = big.NewInt(0)

func isZero(i *big.Int) bool {
	return i.Cmp(BZERO) == 0
}

func BindDefaults(vm *otto.Otto) {
	bindHelpers(vm)
}


func bindHelpers(vm *otto.Otto) {
	vm.Set("Add", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		result, _ := vm.ToValue("0x" + p0.Add(p0, p1).String())
		return result
	})

	vm.Set("Sub", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		p0.Sub(p0, p1)
		if p0.Sign() < 0 {
			otto.NaNValue() // TODO
		}
		result, _ := vm.ToValue("0x" + p0.String())
		return result
	})

	vm.Set("Mul", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		result, _ := vm.ToValue("0x" + p0.Mul(p0, p1).String())
		return result
	})

	vm.Set("Div", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		if isZero(p1) {
			return otto.NaNValue()
		}
		result, _ := vm.ToValue("0x" + p0.Div(p0, p1).String())
		return result
	})

	vm.Set("Mod", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		if isZero(p1) {
			return otto.NaNValue()
		}
		result, _ := vm.ToValue("0x" + p0.Mod(p0, p1).String())
		return result
	})

	vm.Set("Exp", func(call otto.FunctionCall) otto.Value {
		p0, p1, errP := parseBin(call)
		if errP != nil {
			return otto.UndefinedValue()
		}
		result, _ := vm.ToValue("0x" + p0.Exp(p0, p1, nil).String())
		return result
	})

	vm.Set("IsZero", func(call otto.FunctionCall) otto.Value {
		prm, err0 := call.Argument(0).ToString()
		if err0 != nil {
			return otto.UndefinedValue()
		}
		isZero := prm == "0" || prm == "0x" || prm == "0x0"
		result, _ := vm.ToValue(isZero)

		return result
	})

	// Crypto
	vm.Set("SHA3", func(call otto.FunctionCall) otto.Value {
		prm, err0 := call.Argument(0).ToString()
		if err0 != nil {
			return otto.UndefinedValue()
		}
		h, err := hex.DecodeString(prm)
		if err != nil {
			return otto.UndefinedValue()
		}
		d := sha3.NewKeccak256()
		d.Write(h)
		result, _ := vm.ToValue(hex.EncodeToString(d.Sum(nil)))

		return result
	})	
}

func parseUn(call otto.FunctionCall) (*big.Int, error) {
	str, err0 := call.Argument(0).ToString()
	if err0 != nil {
		return nil, err0
	}
	val := atob(str)
	return val, nil
}

func parseBin(call otto.FunctionCall) (*big.Int, *big.Int, error) {
	left, err0 := call.Argument(0).ToString()
	if err0 != nil {
		return nil, nil, err0
	}
	right, err1 := call.Argument(1).ToString()

	if err1 != nil {
		return nil, nil, err1
	}
	p0 := atob(left)
	p1 := atob(right)
	return p0, p1, nil
}

func atob(str string) *big.Int {
	i := new(big.Int)
	fmt.Sscan(str, i)
	return i
}
