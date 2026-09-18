package calculator

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		a, b      float64
		want      float64
		wantErr   bool
	}{
		{"add", "add", 2, 3, 5, false},
		{"subtract", "subtract", 5, 3, 2, false},
		{"multiply", "multiply", 4, 3, 12, false},
		{"divide", "divide", 6, 3, 2, false},
		{"divide by zero", "divide", 1, 0, 0, true},
		{"exponentiate", "exponentiate", 2, 10, 1024, false},
		{"exponentiate negative exponent", "exponentiate", 2, -1, 0.5, false},
		{"exponentiate zero to negative power", "exponentiate", 0, -1, 0, true},
		{"square root", "square_root", 81, 0, 9, false},
		{"square root of negative", "square_root", -4, 0, 0, true},
		{"percentage", "percentage", 20, 50, 10, false},
		{"unsupported", "modulo", 1, 1, 0, true},
		{"exponent too large", "exponentiate", 10, 1000000, 0, true},
		{"exponent at boundary", "exponentiate", 2, 10, 1024, false},
		{"operand a too large", "add", 1e16, 1, 0, true},
		{"operand b too large", "add", 1, -1e16, 0, true},
		{"operand at boundary", "add", 1e15, 0, 1e15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.operation, tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
