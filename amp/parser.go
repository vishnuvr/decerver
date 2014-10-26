package amp

import (
	"encoding/hex"
	"fmt"
	"github.com/eris-ltd/thelonious/monk"
	"math/big"
	"strings"
)

var zeroBytesStr40 = "0000000000000000000000000000000000000000"

const (
	doug     = "doug"
	gendoug  = "gendoug"
	this     = "this"
	param    = "param"
	preCall  = "precall"
	postCall = "postcall"
)

type parseState uint8

const (
	psInit parseState = iota
	psPrecall
	psCall
	psPostcall
	psSuccess
	psResult
)

func PrintPS(ps parseState) string {
	switch ps {
	case psInit:
		return "Init"
	case psPrecall:
		return "Precall"
	case psCall:
		return "Call"
	case psPostcall:
		return "Postcall"
	case psSuccess:
		return "Success"
	case psResult:
		return "Result"
	default:
		return ""
	}
}

var EOFI = item{itemEOF, ""}

func WrapItems(items ...[]item) [][]item {
	arr := make([][]item, 0)
	for _, it := range items {
		arr = append(arr, it)
	}
	return arr
}

func WrapItemsArr(items ...[][]item) [][][]item {
	arr := make([][][]item, 0)
	for _, it := range items {
		arr = append(arr, it)
	}
	return arr
}

// This is a simple combined parser/interpreter at this point. It parses and
// interpretes action model script. Note lexing is done one command at a time,
// because it's an easy way to ensure that both lex + parse errors get the correct
// state + command when reported (meaning if it's done as command nr. 2 in precall,
// that's exactly what will be reported). This is far from an optimal system.
type parser struct {
	ethChain *monk.EthChain
	// Tokens currently being processed
	items []item
	// Position in the token list
	iPos int
	// Constants
	globals map[string]string
	// Parameters
	params map[string]string
	// pre and postcall arrays
	precall  []string
	postcall []string
	// current state
	state parseState
	// current command
	cmd int
}

func NewParser(ethChain *monk.EthChain, globals, params map[string]string) *parser {
	p := &parser{}
	p.ethChain = ethChain
	p.globals = globals
	p.params = params
	return p
}

// Precall or postcall arrays
func (p *parser) ParsePP(prec []string) error {

	if p.state == psInit {
		p.state = psPrecall
		p.precall = []string{}
	} else if p.state == psCall {
		p.state = psPostcall
		p.postcall = []string{}
	} else {
		return p.errorf("Error in parse order.")
	}

	for idx, command := range prec {
		lx := Lex(command)
		items := lx.ItemsFlush()
		last := items[len(items)-1]
		if last.typ == itemError {
			return p.errorf(last.val)
		}

		p.cmd = idx
		// Get left and right side...
		itemsL, itemsR, err := p.split(items)

		if err != nil {
			return err
		}
		// Parse left
		p.setItems(itemsL)
		contract, errL := p.parseNumExpr()
		if errL != nil {
			return errL
		}

		// Parse right
		p.setItems(itemsR)
		address, errR := p.parseArit()
		if errR != nil {
			return errL
		}

		// Query blockchain and append to post state
		result := p.getStorageAt(contract, address)
		if p.state == psPrecall {
			p.precall = append(p.precall, result)
		} else {
			p.postcall = append(p.postcall, result)
		}

	}

	return nil
}

// Call. Returns an array of string-arrays, where each string-array
// is the contract address to call followed by its calldata.
func (p *parser) ParseCall(callData interface{}) ([]string, error) {
	if p.state == psPrecall {
		p.state = psCall
	} else {
		return nil, p.errorf("Error in parse order.")
	}

	var calls [][]string

	switch val := callData.(type) {
	case []string:
		calls = make([][]string,1)
		calls[0] = val
	case [][]string:
		calls = val
	default:
		return nil, p.errorf("Call array is not an array of strings, or string arrays")
	}

	// For each call ([]string)
	if len(calls) > 0 {
		for idx, call := range calls {
			if len(call) == 0 {
				continue
			}
			p.cmd = idx
			// Start with call[0], which is the contract address
			strs := []string{}
			lx := Lex(call[0])
			itms := lx.ItemsFlush()
			last := itms[len(itms)-1]
			if last.typ == itemError {
				return nil, p.errorf(last.val)
			}

			p.setItems(itms)
			str, err := p.parseNumExpr()
			if err != nil {
				return nil, err
			}
			strs = append(strs, str)

			if len(call) > 1 {
				for _, c := range call[1:] {

					lx := Lex(c)
					is := lx.ItemsFlush()
					last := is[len(is)-1]
					if last.typ == itemError {
						return nil, p.errorf(last.val)
					}

					p.setItems(is)
					str, err = p.parseNumExpr()
					if err != nil {
						return nil, err
					}
					strs = append(strs, str)
				}
				p.ethChain.Msg(strs[0], strs[1:])
			} else {
				// If no args, just pass an empty array.
				p.ethChain.Msg(strs[0], []string{})
			}
		}
	}
	return nil, nil
}

// Parse success
func (p *parser) ParseSuccess(expr string) (bool, error) {
	p.cmd = 0

	if p.state == psPostcall {
		p.state = psSuccess
	} else {
		return false, p.errorf("Error in parse order.")
	}

	lx := Lex(expr)
	items := lx.ItemsFlush()
	last := items[len(items)-1]
	if last.typ == itemError {
		return false, p.errorf(last.val)
	}

	var left, right string
	var err error

	p.setItems(items)

	// First param.
	it := p.current()

	left, err = p.parseNumExpr()

	if err != nil {
		return false, err
	}

	it = p.next()
	if !isCond(it) {
		return false, p.errorf("Success needs a conditional statement.")
	}
	cnd := it.typ

	// Second param.
	it = p.next()
	right, err = p.parseNumExpr()
	if err != nil {
		return false, err
	}

	// If we made it here we're ready to compare.
	p1 := p.strToBigInt(left)
	p2 := p.strToBigInt(right)

	cmp := p1.Cmp(p2)

	switch cnd {
	case itemEq:
		return cmp == 0, nil
	case itemNeq:
		return cmp != 0, nil
	case itemGeq:
		return cmp >= 0, nil
	case itemGt:
		return cmp > 0, nil
	case itemLeq:
		return cmp <= 0, nil
	case itemLt:
		return cmp < 0, nil
	}
	// Should not come to here.
	return false, nil
}

// Result
func (p *parser) ParseResult(rs []string) ([]string, error) {
	if p.state == psSuccess {
		p.state = psResult
	} else {
		return nil, p.errorf("Error in parse order.")
	}
	result := []string{}
	for idx, r := range rs {
		p.cmd = idx
		lx := Lex(r)
		items := lx.ItemsFlush()
		last := items[len(items)-1]
		if last.typ == itemError {
			return nil, p.errorf(last.val)
		}
		p.setItems(items)
		res, err := p.parseNumExpr()
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

func (p *parser) parseNumExpr() (string, error) {
	it := p.current()
	if !isNumberType(it) {
		return "", p.errorf("Left side is not (and does not evaluate to) a number")
	}
	if it.typ == itemIdentifier {
		return p.parseIdentifier()
	} else {
		if it.typ == itemString {
			return p.parseString(it)
		}
		return p.parseInt(it)
	}
}

func (p *parser) parseArit() (string, error) {
	it := p.current()
	if !isNumberType(it) {
		return "", p.errorf("Right side must start with either an identifier, a string, or a number")
	}
	var res string
	if it.typ == itemIdentifier {
		var errPI error
		res, errPI = p.parseIdentifier()
		if errPI != nil {
			return "", errPI
		}
	} else if it.typ == itemNumber {
		var errPN error
		res, errPN = p.parseInt(it)
		if errPN != nil {
			return "", errPN
		}
	} else {
		var errPS error
		res, errPS = p.parseString(it)
		if errPS != nil {
			return "", errPS
		}
	}

	//First step done. Now check if we're done, or if there is a + or - here.
	it = p.next()
	if it.typ == itemEOF {
		return res, nil
	}

	if it.typ == itemPlus || it.typ == itemMinus {
		p1 := p.strToBigInt(res)
		op := it.typ
		it = p.next()
		res2, errR2 := p.parseInt(it)
		if errR2 != nil {
			return "", errR2
		}
		p2 := p.strToBigInt(res2)
		if op == itemPlus {
			p1.Add(p1, p2)
		} else {
			p1.Sub(p1, p2)
			if p1.Sign() == -1 {
				return "", p.errorf("Right side is a negative number.")
			}
		}
		res = "0x" + hex.EncodeToString(p1.Bytes())
		return res, nil
	}

	return "", p.errorf("Right side does not evaluate to a number.")
}

// Big one, how to handle all the $-stuff
func (p *parser) parseIdentifier() (string, error) {
	ctVal := p.current().val
	switch ctVal {
	case doug:
		fallthrough
	case gendoug:
		fallthrough
	case this:
		if val, ok := p.globals[ctVal]; ok {
			return val, nil
		} else {
			return "", p.errorf(fmt.Sprintf("Constant does not exist: %s", ctVal))
		}
	case preCall:
		// We will have to check if precall[n] is done, which is true
		// if it exists.
		lb := p.next()
		if lb.typ != itemLb {
			return "", p.errorf("Expected [.")
		}
		iVal := p.next()
		rb := p.next()
		if rb.typ != itemRb {
			return "", p.errorf("Expected ].")
		}
		if iVal.typ != itemNumber {
			return "", p.errorf("Expected number inside brackets.")
		}
		idx := int(p.strToBigInt(iVal.val).Uint64())
		if idx >= len(p.precall) {
			return "", p.errorf("Precall index out of range.")
		}
		// Success.
		return p.precall[idx], nil
	case postCall:
		if p.state < psPostcall {
			return "", p.errorf("Referencing postcall before the post call step.")
		}
		// We will have to check if postcall[n] is done, which is true
		// if it exists.
		lb := p.next()
		if lb.typ != itemLb {
			return "", p.errorf("Expected [.")
		}
		iVal := p.next()
		rb := p.next()
		if rb.typ != itemRb {
			return "", p.errorf("Expected ].")
		}
		if iVal.typ != itemNumber {
			return "", p.errorf("Expected number inside brackets.")
		}
		idx := int(p.strToBigInt(iVal.val).Uint64())
		if idx >= len(p.precall) {
			return "", p.errorf("Postcall index out of range.")
		}
		// Success.
		return p.postcall[idx], nil
	case param:
		// We will have to check if postcall[n] is done, which is true
		// if it exists.
		lb := p.next()
		if lb.typ != itemLb {
			return "", p.errorf("Expected [.")
		}
		sVal := p.next()
		rb := p.next()
		if rb.typ != itemRb {
			return "", p.errorf("Expected ].")
		}
		if sVal.typ != itemString {
			return "", p.errorf("Expected quoted string inside brackets.")
		}
		// Success.
		par := p.params[sVal.val]
		if val, ok := p.params[par]; ok {
			return val, nil
		} else {
			return "", p.errorf(fmt.Sprintf("Parameter does not exist: %s", sVal.val))
		}
	}
	return "", nil
}

// Parses a sequence of alphanumeric characters and turn them into hex.
func (p *parser) parseString(it item) (string, error) {
	if it.typ != itemString {
		return "", p.errorf(fmt.Sprintf("Not a valid string."))
	}
	hexStr := hex.EncodeToString([]byte(it.val))
	hexStr = "0x" + hexStr + zeroBytesStr40[:40-len(hexStr)]

	return hexStr, nil
}

func (p *parser) parseInt(it item) (string, error) {
	if it.typ != itemNumber {
		return "", p.errorf(fmt.Sprintf("Not a valid number."))
	}
	return it.val, nil
}

func (p *parser) parseInt160(it item) (string, error) {
	if it.val[0:2] != "0x" || len(it.val) != 42 {
		return "", p.errorf(fmt.Sprintf("Not a valid contract address."))
	}
	return it.val, nil
}

func (p *parser) setItems(is []item) {
	p.items = is
	p.iPos = 0
}

func (p *parser) peek() item {
	it := p.next()
	p.backup()
	return it
}

func (p *parser) current() item {
	return p.items[p.iPos]
}

func (p *parser) next() item {
	if p.iPos == len(p.items)-1 {
		return EOFI
	}
	p.iPos += 1
	return p.items[p.iPos]
}

func (p *parser) backup() {
	p.iPos -= 1
}

func (p *parser) indexOf(it itemType, items []item) int {
	idxOf := -1
	for i, itm := range items {
		if itm.typ == it {
			idxOf = i
			break
		}
	}
	return idxOf
}

func (p *parser) split(items []item) ([]item, []item, error) {

	idxOf := p.indexOf(itemCol, items)
	// Commands must have at least one colon in them.
	if idxOf == -1 {
		return nil, nil, p.errorf(fmt.Sprintf("Missing colon (:) separator."))
	}
	// Colons must have at least one token on each side.
	if idxOf == 0 {
		return nil, nil, p.errorf(fmt.Sprintf("Misplaced colon (:) separator: %d", idxOf))
	}
	items0 := items[0:idxOf]
	items1 := items[idxOf+1:]
	// Can't be more then one colon.
	if p.indexOf(itemCol, items1) != -1 {
		return nil, nil, p.errorf(fmt.Sprintf("More then one colon (:) separator: %d", idxOf))
	}
	return items0, items1, nil
}

func (p *parser) getStorageAt(s, a string) string {
	val := p.ethChain.GetStorageAt(s, a)
	if len(val) > 2 && strings.HasPrefix(val,"0x") {
		return "0x" + val
	} else if val == "0x" {
		return "0x0"
	} else {
		return "0x" + val
	}
}

func (p *parser) strToBigInt(str string) *big.Int {
	i := new(big.Int)
	fmt.Sscan(str, i)
	return i
}

func (p *parser) errorf(s string) error {
	return fmt.Errorf("Parser error (State: %s, Command %d): %s", PrintPS(p.state), p.cmd, s)
}

// These all evaluate to a hex number.
func isNumberType(it item) bool {
	if it.typ == itemIdentifier || it.typ == itemNumber || it.typ == itemString {
		return true
	}
	return false
}
