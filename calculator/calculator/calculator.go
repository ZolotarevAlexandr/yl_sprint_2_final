package calculator

func Calculate(expression string) (float64, error) {
	stringTokens, err := Tokenize(expression)
	if err != nil {
		return 0, err
	}

	tokens, err := ShuntingYard(stringTokens)
	if err != nil {
		return 0, err
	}

	return Evaluate(tokens)
}
