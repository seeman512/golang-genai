package converter

import (
	"math"
	"testing"
)

// TestConvertCurrency — таблично-керований тест, який запускається
// автоматично у GitHub Actions при кожному push і Pull Request.
func TestConvertCurrency(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		rate    float64
		want    float64
		wantErr bool
	}{
		{
			name:   "звичайна конвертація",
			amount: 100,
			rate:   0.91,
			want:   91,
		},
		{
			name:   "нульова сума — не помилка",
			amount: 0,
			rate:   0.91,
			want:   0,
		},
		{
			name:    "від'ємна сума — помилка",
			amount:  -50,
			rate:    0.91,
			wantErr: true,
		},
		{
			name:    "нульовий курс — помилка",
			amount:  100,
			rate:    0,
			wantErr: true,
		},
		{
			name:    "від'ємний курс — помилка",
			amount:  100,
			rate:    -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertCurrency(tt.amount, tt.rate)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ConvertCurrency(%v, %v) = %v, <nil>; хотіли помилку, а не отримали",
						tt.amount, tt.rate, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("ConvertCurrency(%v, %v) повернула неочікувану помилку: %v",
					tt.amount, tt.rate, err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("ConvertCurrency(%v, %v) = %v, want %v",
					tt.amount, tt.rate, got, tt.want)
			}
		})
	}
}
