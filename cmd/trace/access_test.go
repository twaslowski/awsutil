package trace

import (
	"testing"
	"time"
)

func Test_parseStartDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Should parse timestamp via dateparse",
			input:   "2025-02-02",
			want:    time.Date(2025, 02, 02, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Should parse timestamp via custom timedelta",
			input:   "7d",
			want:    time.Now().Add(-time.Hour * 24 * 7),
			wantErr: false,
		},
		{
			name:    "Should raise error for empty input",
			input:   "",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Should raise error for invalid input",
			input:   "abc",
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStartDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStartDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Day() != tt.want.Day() {
				t.Errorf("parseStartDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
