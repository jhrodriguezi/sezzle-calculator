package calculator

import (
	"fmt"
	"math"
)

// Operation identifies a supported arithmetic operation.
type Operation string

const (
	Add          Operation = "add"
	Subtract     Operation = "subtract"
	Multiply     Operation = "multiply"
	Divide       Operation = "divide"
	Exponentiate Operation = "exponentiate"
	SquareRoot   Operation = "square_root"
	Percentage   Operation = "percentage"
)

// fn is the signature every operation implements. Unary operations (e.g.
// SquareRoot) simply ignore b.
type fn func(a, b float64) (float64, error)

// operations maps each supported Operation to its implementation, making it
// straightforward to add new operations without touching the dispatch logic.
var operations = map[Operation]fn{
	Add:          add,
	Subtract:     subtract,
	Multiply:     multiply,
	Divide:       divide,
	Exponentiate: exponentiate,
	SquareRoot:   squareRoot,
	Percentage:   percentage,
}

// Evaluate performs the given arithmetic operation on a and b.
// For unary operations (SquareRoot), b is ignored.
func Evaluate(operation string, a, b float64) (float64, error) {
	op, ok := operations[Operation(operation)]
	if !ok {
		return 0, fmt.Errorf("unsupported operation: %s", operation)
	}

	result, err := op(a, b)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, fmt.Errorf("operation produced an invalid result")
	}

	return result, nil
}

func add(a, b float64) (float64, error) {
	return a + b, nil
}

func subtract(a, b float64) (float64, error) {
	return a - b, nil
}

func multiply(a, b float64) (float64, error) {
	return a * b, nil
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// exponentiate computes a raised to the power of b (a^b).
func exponentiate(a, b float64) (float64, error) {
	if a == 0 && b < 0 {
		return 0, fmt.Errorf("cannot raise zero to a negative power")
	}
	return math.Pow(a, b), nil
}

// squareRoot computes the square root of a. b is ignored.
func squareRoot(a, _ float64) (float64, error) {
	if a < 0 {
		return 0, fmt.Errorf("cannot take the square root of a negative number")
	}
	return math.Sqrt(a), nil
}

// percentage computes a as a percentage of b (a% of b).
func percentage(a, b float64) (float64, error) {
	return (a / 100) * b, nil
}
