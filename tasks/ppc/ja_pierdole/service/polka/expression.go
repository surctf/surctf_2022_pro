package polka

import (
	"strconv"
	"strings"
)

const (
	PostfixType = 0
	PrefixType  = 1
)

// Expression - used to store expression in string format and type of expression.
type Expression struct {
	String string
	Type   int
}

// Result - calculate result of Expression.
func (e *Expression) Result() int {
	if e.Type == PrefixType {
		return calc(strings.Split(e.String, " "), e.Type)
	}

	return calc(reverse(strings.Split(e.String, " ")), e.Type)
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func execute(op string, n1 int, n2 int) int {
	switch op {
	case "+":
		return n1 + n2
	case "-":
		return n1 - n2
	case "*":
		return n1 * n2
	}
	return 0
}

func isArithmetic(c string) bool {
	switch c {
	case "+", "-", "*":
		return true
	}
	return false
}

func calc(stack []string, t int) int {
	var opInd, n1, n2 int
	var err error

	for {
		for i, c := range stack {
			if isArithmetic(c) {
				n1, err = strconv.Atoi(stack[i+1])
				if err != nil {
					continue
				}

				n2, err = strconv.Atoi(stack[i+2])
				if err != nil {
					continue
				}

				opInd = i
				break
			}
		}

		var result int
		if t == PrefixType {
			result = execute(stack[opInd], n1, n2)
		} else {
			result = execute(stack[opInd], n2, n1)
		}

		if len(stack) == 3 {
			return result
		}

		stack = append(append(stack[0:opInd], strconv.Itoa(result)), stack[opInd+3:]...)
		opInd = -1
	}

}
