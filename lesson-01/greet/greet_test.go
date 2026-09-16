package greet

import "testing"

// TestGreet — таблично-керований (table-driven) тест, який запускається
// автоматично у GitHub Actions при кожному push і Pull Request.
func TestGreet(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "звичайне ім'я",
			input: "Alice",
			want:  "Hello, Alice! Welcome to Go.",
		},
		{
			name:  "інше ім'я",
			input: "Bohdan",
			want:  "Hello, Bohdan! Welcome to Go.",
		},
		{
			name:  "порожнє ім'я",
			input: "",
			want:  "Hello, stranger! Welcome to Go.",
		},
		{
			name:  "пробіли навколо імені обрізаються",
			input: "  Bob  ",
			want:  "Hello, Bob! Welcome to Go.",
		},
		{
			name:  "рядок лише з пробілів прирівнюється до порожнього",
			input: "   ",
			want:  "Hello, stranger! Welcome to Go.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Greet(tt.input)
			if got != tt.want {
				t.Errorf("Greet(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
