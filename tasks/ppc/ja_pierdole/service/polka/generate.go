package polka

import (
	"math/rand"
	"strconv"
	"strings"
)

const (
	MAX_INT            = 25
	COMPLEX_EXP_CHANCE = 50
	SYMBOLS            = "+-*"
)

// Seed - sets seed for rand.
func Seed(n int64) {
	rand.Seed(n)
}

// generateSimplePrefixExpr - generates simple prefix expression.
func generateSimplePrefixExpr() string {
	op := SYMBOLS[rand.Int()%len(SYMBOLS)]
	n1 := rand.Int() % MAX_INT
	n2 := rand.Int() % MAX_INT

	exp := string(op) + " " + strconv.Itoa(n1) + " " + strconv.Itoa(n2)
	return exp
}

// generatePrefixExpr - generates simple/complex prefix expression.
func generatePrefixExpr(complexExpChance int) string {
	isComplexExp1 := (rand.Int() % 100) < complexExpChance
	isComplexExp2 := (rand.Int() % 100) < complexExpChance

	var exp1 string
	if isComplexExp1 {
		exp1 = generatePrefixExpr(complexExpChance / 2)
	} else {
		exp1 = generateSimplePrefixExpr()
	}

	var exp2 string
	if isComplexExp2 {
		exp2 = generatePrefixExpr(complexExpChance / 2)
	} else {
		exp2 = generateSimplePrefixExpr()
	}

	op := SYMBOLS[rand.Int()%len(SYMBOLS)]
	return string(op) + " " + exp1 + " " + exp2
}

// GenerateExpr - generates simple/complex prefix/postfix expression.
func GenerateExpr() Expression {
	s := generatePrefixExpr(COMPLEX_EXP_CHANCE)

	if isPrefix := rand.Int() % 2; isPrefix%2 == 0 {
		return Expression{String: s, Type: PrefixType}
	}

	s = strings.Join(reverse(strings.Split(s, " ")), " ")
	return Expression{String: s, Type: PostfixType}
}
