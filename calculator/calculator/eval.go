package calculator

import "errors"

func Evaluate(tokens []Token) (float64, error) {
	stack := make([]float64, 0)
	for _, token := range tokens {
		switch {
		case token.IsOperand:
			val, _ := token.GetOperand()
			stack = append(stack, val)
		case token.IsOperator:
			if len(stack) < 2 {
				return 0.0, errors.New("not enough operands")
			}
			num1, num2 := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			operator, err := token.getOperator()
			if err != nil {
				return 0.0, err
			}
			val, err := EvaluateOperation(operator, num1, num2)
			if err != nil {
				return 0.0, err
			}
			stack = append(stack, val)
		}
	}
	if len(stack) != 1 {
		return 0.0, errors.New("invalid stack")
	}
	return stack[len(stack)-1], nil
}

func EvaluateOperation(operator string, num1, num2 float64) (float64, error) {
	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0.0, errors.New("division by zero")
		}
		return num1 / num2, nil
	}
	return 0.0, errors.New("invalid operand")
}
