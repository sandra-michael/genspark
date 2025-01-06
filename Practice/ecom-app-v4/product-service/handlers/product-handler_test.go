package handlers

import (
	"testing"
)

func TestRupeesToPaise(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      uint64
		expectErr bool
	}{
		{
			name:      "Valid price with no decimal",
			input:     "100",
			want:      10000,
			expectErr: false,
		},
		{
			name:      "Valid price with two decimals",
			input:     "100.25",
			want:      10025,
			expectErr: false,
		},
		{
			name:      "Valid price with one decimal",
			input:     "99.2",
			want:      9920,
			expectErr: false,
		},
		{
			name:      "Valid price with leading and trailing spaces",
			input:     "  50.50  ",
			want:      5050,
			expectErr: false,
		},
		{
			name:      "Invalid price with letters",
			input:     "10a.50",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Invalid price with too many decimals",
			input:     "10.123",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Invalid price with multiple dots",
			input:     "10.50.25",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Invalid price with empty string",
			input:     "",
			want:      0,
			expectErr: true,
		},
		{
			name:      "Valid price with zero paisa",
			input:     "50.00",
			want:      5000,
			expectErr: false,
		},
		{
			name:      "Invalid price with negative value",
			input:     "-50.25",
			want:      0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RupeesToPaise(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("RupeesToPaise() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("RupeesToPaise() = %v, want %v", got, tt.want)
			}
		})
	}
}
