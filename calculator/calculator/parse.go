package calculator

import "fmt"

func ShuntingYard(tokens []Token) ([]Token, error) {
	outputStack := make([]Token, 0)
	operatorsStack := make([]Token, 0)

	for _, token := range tokens {
		switch {
		case token.IsOperator:
			for len(operatorsStack) > 0 {
				top := operatorsStack[len(operatorsStack)-1]
				if top.IsOperator && (token.Priority <= top.Priority) {
					outputStack = append(outputStack, top)
					operatorsStack = operatorsStack[:len(operatorsStack)-1]
				} else {
					break
				}
			}
			operatorsStack = append(operatorsStack, token)

		case token.isOpeningBracket():
			operatorsStack = append(operatorsStack, token)

		case token.isClosingBracket():
			for len(operatorsStack) > 0 && !operatorsStack[len(operatorsStack)-1].isOpeningBracket() {
				outputStack = append(outputStack, operatorsStack[len(operatorsStack)-1])
				operatorsStack = operatorsStack[:len(operatorsStack)-1]
			}
			// Remove the open parenthesis
			if len(operatorsStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			operatorsStack = operatorsStack[:len(operatorsStack)-1]

		default:
			outputStack = append(outputStack, token)
		}
	}

	// Pop any remaining operators in the stack
	for len(operatorsStack) > 0 {
		top := operatorsStack[len(operatorsStack)-1]
		if top.IsBracket {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		outputStack = append(outputStack, top)
		operatorsStack = operatorsStack[:len(operatorsStack)-1]
	}

	return outputStack, nil
}
