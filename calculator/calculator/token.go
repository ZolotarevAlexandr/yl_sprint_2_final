package calculator

import (
	"errors"
	"strconv"
	"strings"
	"text/scanner"
)

type Token struct {
	IsOperator bool
	IsOperand  bool
	IsBracket  bool
	Priority   int
	Value      any
}

func (t Token) getOperator() (string, error) {
	if t.IsOperator == false {
		return "", errors.New("token is not an operator")
	}
	return t.Value.(string), nil
}

func (t Token) GetOperand() (float64, error) {
	if t.IsOperand == false {
		return 0, errors.New("token is not an operand")
	}
	return t.Value.(float64), nil
}

func (t Token) getBracket() (string, error) {
	if t.IsBracket == false {
		return "", errors.New("token is not a bracket")
	}
	return t.Value.(string), nil
}

func (t Token) isOpeningBracket() bool {
	val, err := t.getBracket()
	if err != nil {
		return false
	}
	return val == "("
}

func (t Token) isClosingBracket() bool {
	val, err := t.getBracket()
	if err != nil {
		return false
	}
	return val == ")"
}

var priorities = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func strToToken(str string) (Token, error) {
	if priority, ok := priorities[str]; ok {
		return Token{IsOperator: true, Priority: priority, Value: str}, nil
	}
	if str == "(" || str == ")" {
		return Token{IsBracket: true, Value: str}, nil
	}
	num, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return Token{IsOperand: true, Value: num}, nil
	}
	return Token{}, errors.New("unsupported token value")
}

func Tokenize(str string) ([]Token, error) {
	result := make([]Token, 0)
	var scan scanner.Scanner
	var token rune

	scan.Init(strings.NewReader(str))

	for token != scanner.EOF {
		token = scan.Scan()
		val := strings.TrimSpace(scan.TokenText())
		if len(val) <= 0 {
			continue
		}
		tok, err := strToToken(val)
		if err != nil {
			return []Token{}, err
		}
		result = append(result, tok)
	}

	return result, nil
}
