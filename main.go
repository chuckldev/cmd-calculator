package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// initialize regex patterns
var rxParens = regexp.MustCompile(`(\()|(\))`)
var rxOps = regexp.MustCompile(`[\-\+\/\*]`)

// setup operations function map
// this map will map a function ( i.e. 'Add', 'Subt', 'Mult', 'Div' ) to an operator ( i.e. '*', '+', '-', '/' )
var opm map[string]func(float64,float64)float64
var opPrecendence []string

//initialize operations map and operator precedence
func init() {
	opm = make(map[string]func(float64,float64)float64)
	opm[`*`] = Mult
	opm[`/`] = Div
	opm[`+`] = Add
	opm[`-`] = Subt

	opPrecendence = []string{`*`,`/`,`+`,`-`}
}

func CheckErr( e error ) {
	if e != nil {
		panic(e)
	}
}

func Add( op1, op2 float64) float64 {
	return op1 + op2
}

func Subt( op1, op2 float64 ) float64 {
	return op1 - op2
}

func Mult( op1, op2 float64 ) float64 {
	return op1 * op2
}

func Div( op1, op2 float64 ) float64 {
	var res float64
	if op2 == 0 {
		panic( "Cannot divide by Zero ( 0 )" )
	}
	res = op1 / op2
	return res
}

// check for parentheses
func HasParens( s string ) bool {
	return rxParens.MatchString( s )
}

// parse the indices of parentheses
func ParensIndex(s string) []int {
	var op []int
	parsed := false
	re := regexp.MustCompile(`\(|\)`)
	parens := re.FindAllStringIndex(s, -1)
	res := make([]int, 2)
	for i := range parens {
		ss := s[parens[i][0]:parens[i][1]]
		switch ss {
		case `(`:
			op = append(op, parens[i][1])
		case `)`:
			oparen := op[len(op)-1]
			if len(op) > 0 {
				op = op[:len(op)-1]
			}
			res[0] = oparen
			res[1] = parens[i][0]
			parsed = true
		}
		if parsed {
			break
		}
	}
	return res
}

// parse the expressions from the parentheses
func ParseExpressions(s string) [][]string {
	parens := ParensIndex(s)
	var res [][]string
	expr := make([]string, 2)
	expr[0] = s[parens[0]-1:parens[1]+1]
	expr[1] = s[parens[0]:parens[1]]
	res = append(res, expr)
	return res
}

// check an array if it contains the op value
func Contains(ops []string, op string) int {
	for i := range ops {
		if ops[i] == op {
			return i
		}
	}

	return -1
}

// parse the operands fromt he expressions
func ParseOperands(operands []string, idx int) (op1, op2 float64, ok bool) {
	op1, err := strconv.ParseFloat( operands[idx], 64 )
	if err != nil {
		return math.NaN(), math.NaN(), false
	}
	op2, err = strconv.ParseFloat( operands[idx+1], 64 )
	if err != nil {
		return math.NaN(), math.NaN(), false
	}

	return op1, op2, true
}

// evaluate the expressions
func EvaluateExpression(s string) float64 {
	var res float64
	for {
		operators := rxOps.FindAllString(s, -1)
		if operators == nil {
			break
		}
		operands := rxOps.Split(s, -1)

		for i := 0; i < len(opPrecendence); i++ {
			oper := opPrecendence[i]
			idx := Contains(operators, oper)
			if idx > -1 {
				op1, op2, ok := ParseOperands(operands, idx)
				if !ok {
					err := fmt.Errorf("Something went wrong with parsing operands[%s]", strings.Join(operands, ","))
					panic(err)
				}
				res := opm[oper](op1, op2)
				oldExpr := fmt.Sprintf("%s%s%s", operands[idx], oper, operands[idx+1])
				newExpr := strconv.FormatFloat( res, 'f', -1, 64 )
				s = strings.Replace(s, oldExpr, newExpr, -1)
				break
			}
		}
	}
	res, err := strconv.ParseFloat( s, 64 )
	if err != nil {
		panic(err)
	}

	return res
}

// replace the expression within parentheses, in the original expression,
// with the solution of the parentheses expression
func ReplaceParens(s *string, oldstr, newstr string) {
	*s = strings.Replace(*s, oldstr, newstr, -1)
}

func main() {
	if len(os.Args) < 1 {
		panic( "Calculator requires at least one argument ( string )." )
	}

	// the problem inputted from cmd line
	prob := os.Args[1]

	// check to see if prob has any parentheses
	// continue looping until all expressions in parentheses
	// have been evaluated
	for HasParens(prob) {
		expr := ParseExpressions(prob)
		for i := range expr {
			res := EvaluateExpression(expr[i][1])
			ReplaceParens(&prob, expr[i][0], strconv.FormatFloat( res, 'f', -1, 64 ))
		}
	}

	// evaluate the problem entered
	solution := EvaluateExpression(prob)

	// print the solution
	fmt.Println(os.Args[1], " = ", solution)
}