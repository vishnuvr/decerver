package amp

import (
	"fmt"
	"strings"
)

var DEBUG bool = true

const (
	itemError      itemType = iota
	itemCol                 // :
	itemLb                  // [
	itemRb                  // ]
	itemPlus                // +
	itemMinus               // -
	itemConst				// ?
	itemIdentifier          // identifier
	itemNumber              // number
	itemString				// non quoted string (a sequence of letters and _'s).
	itemDQString            // double quoted string: "likethis"
	itemEq                  // ==
	itemNeq                 // !=
	itemLt                  // <
	itemGt                  // >
	itemLeq                 // <=
	itemGeq                 // >=
	itemEOF                 // -1
	// We don't tokenize '$'
)

func PrintIT(i itemType) string {
	switch i {
	case itemError:
		return "Error"
	case itemCol:
		return "Colon"
	case itemLb:
		return "Left Bracket"
	case itemRb:
		return "Right Bracket"
	case itemPlus:
		return "Plus"
	case itemMinus:
		return "Minus"
	case itemConst:
		return "Const Sign"
	case itemIdentifier:
		return "Identifier"
	case itemNumber:
		return "Number"
	case itemDQString:
		return "Quoted string"
	case itemString:
		return "String"
	case itemEq:
		return "Equal to"
	case itemNeq:
		return "Not equal to"
	case itemLt:
		return "Less then"
	case itemGt:
		return "Greater then"
	case itemLeq:
		return "Less then or equals"
	case itemGeq:
		return "Greater then or equals"
	case itemEOF:
		return "EOF"
	default:
		return "Not an itemtype"
	}
}

const (
	COMMA = ','
	COLON = ':'
	PLUS  = '+'
	MINUS = '-'
	QUO   = '"'
	VAR   = '$'
	LT    = '<'
	GT    = '>'
	EQ    = '='
	NOT   = '!'
	LB    = '['
	RB    = ']'
	EOF   = 0
)

const digits = "0123456789"

const digitsHex = "0123456789abcdefABCDEF"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const id = letters + "_"

type itemType uint8

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	return "Type: " + PrintIT(i.typ) + ", Value: " + i.val;
}

type tokenStream interface {
	Next() (item, error)
	Add(it item)
	Flush() []item
}

type tokenStreamArray struct {
	items []item
}

func NewTokenStreamArray() *tokenStreamArray {
	return &tokenStreamArray{[]item{}}
}

func (tsa *tokenStreamArray) Next() (item, error) {
	// This should never happen if used correctly.
	if len(tsa.items) == 0 {
		return item{}, fmt.Errorf("Lexer Error: Error reading token stream: null token")
	}
	it := tsa.items[0]
	if len(tsa.items) > 1 {
		tsa.items = tsa.items[1:]
	} else {
		tsa.items = []item{}
	}
	return it, nil
}

func (tsa *tokenStreamArray) Add(it item) {
	tsa.items = append(tsa.items, it)
}

func (tsa *tokenStreamArray) Flush() []item {
	itms := tsa.items
	tsa.items = nil
	if DEBUG {
		fmt.Println("Lexer flushing:")
		for idx, itm := range itms {
			fmt.Printf("%d: %s\n", idx, itm.String())
		}
	}
	return itms
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input string // the string being scanned.
	start int    // token start pos.
	pos   int    // current position in the input.
	items tokenStream
}

func Lex(input string) *lexer {
	l := &lexer{
		input: input + string(0),
		items: &tokenStreamArray{},
	}
	l.run() // Not concurrent
	return l
}

func (l *lexer) run() {
	for state := lexCMD; state != nil; {
		state = state(l)
	}
}

// next returns the next byte in the input.
func (l *lexer) next() byte {
	if l.pos >= len(l.input) {
		return EOF
	}
	b := l.input[l.pos]
	l.pos += 1
	return b
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back one byte.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= 1
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() byte {
	b := l.next()
	l.backup()
	return b
}

func (l *lexer) emit(t itemType) {
	l.items.Add(item{t, l.input[l.start:l.pos]})
	l.start = l.pos
}

// accept consumes the next byte
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexByte(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of bytes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexByte(valid, l.next()) >= 0 {
	}
	l.backup()
}

func lexNumber(l *lexer) stateFn {
	// Is it hex?
	dg := digits
	if l.accept("0") && l.accept("xX") {
		dg = digitsHex
	}
	l.acceptRun(dg)
	l.emit(itemNumber)
	return lexCMD
}

func lexQuote(l *lexer) stateFn {
	l.ignore() // Cut out "
	Loop:
	for {
		switch l.next() {
		case EOF, '\n':
			return l.errorf("Lexer Error: Unterminated quoted string.")
		case QUO:
			break Loop
		}
	}
	// Cut out "
	l.backup()
	l.emit(itemString)
	l.next()
	l.ignore()
	return lexCMD
}

func lexIdent(l *lexer) stateFn {
	l.acceptRun(id)
	l.emit(itemIdentifier)
	return lexCMD
}

func lexString(l *lexer) stateFn {
	l.acceptRun(id)
	l.emit(itemString)
	return lexCMD
}

func lexCond(l *lexer) stateFn {
	// Very small space, just do manually
	switch r := l.next(); {
	case r == EQ:
		if l.next() == EQ {
			l.emit(itemEq)
			return lexCMD(l)
		}
		return l.errorf("Lexer Error: Invalid token: =")
	case r == NOT:
		if l.next() == EQ {
			l.emit(itemNeq)
			return lexCMD(l)
		}
		return l.errorf("Lexer Error: Invalid token: !")
	case r == LT:
		if l.next() == EQ {
			l.emit(itemLeq)
			return lexCMD(l)
		}
		l.backup()
		l.emit(itemLt)
		return lexCMD(l)
	case r == GT:
		if l.next() == EQ {
			l.emit(itemGeq)
			return lexCMD(l)
		}
		l.backup()
		l.emit(itemGt)
	}
	return lexCMD
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items.Add(item{
		itemError,
		fmt.Sprintf(format, args...),
	})
	return nil
}

func (l *lexer) ItemsNext() (item, error) {
	return l.items.Next()
}

func (l *lexer) ItemsFlush() []item {
	return l.items.Flush()
}

func lexCMD(l *lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == EOF:
			return nil
		case isSpace(r):
			l.ignore()
		case r == COLON:
			l.emit(itemCol)
		case r == PLUS:
			l.emit(itemPlus)
		case r == MINUS:
			l.emit(itemMinus)
		case r == LB:
			l.emit(itemLb)
		case r == RB:
			l.emit(itemRb)
		case r == VAR:
			// Consume the $
			l.ignore()
			return lexIdent(l)
		case '0' <= r && r <= '9':
			l.backup()
			return lexNumber(l)
		case r == EQ || r == NOT || r == GT || r == LT:
			l.backup()
			return lexCond(l)
		case isId(r):
			return lexString(l)
		default:
			return l.errorf("Lexer Error: Unexpected character: %c\n", r)
		}
	}
}

// Should do tables but see first where this ends up.

func isSpace(r byte) bool {
	return r == ' ' || r == '\t'
}

func isId(r byte) bool {
	if strings.IndexByte(id, r) >= 0 {
		return true
	}
	return false
}

func isCond(it item) bool {
	switch it.typ {
	case itemEq:
		fallthrough
	case itemNeq:
		fallthrough
	case itemGeq:
		fallthrough
	case itemGt:
		fallthrough
	case itemLeq:
		fallthrough
	case itemLt:
		return true
	default:
		return false
	}
	return false
}
