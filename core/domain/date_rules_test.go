package domain

import (
	"testing"
	"time"
)

func TestNextOccurrence(t *testing.T) {
	tests := []struct {
		name     string
		original time.Time
		from     time.Time
		yearly   bool
		want     time.Time
	}{
		{
			name:     "non yearly returns original",
			original: time.Date(2024, 5, 1, 13, 0, 0, 0, time.UTC),
			from:     time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			yearly:   false,
			want:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "same year upcoming",
			original: time.Date(2020, 12, 5, 0, 0, 0, 0, time.UTC),
			from:     time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			yearly:   true,
			want:     time.Date(2026, 12, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "same day is current occurrence",
			original: time.Date(2020, 6, 11, 18, 0, 0, 0, time.UTC),
			from:     time.Date(2026, 6, 11, 9, 0, 0, 0, time.UTC),
			yearly:   true,
			want:     time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "past this year rolls to next year",
			original: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
			from:     time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			yearly:   true,
			want:     time.Date(2027, 1, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "feb 29 projects to feb 28 in non leap year",
			original: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			from:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			yearly:   true,
			want:     time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "feb 29 remains feb 29 in leap year",
			original: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			from:     time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC),
			yearly:   true,
			want:     time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextOccurrence(tt.original, tt.from, tt.yearly)
			if !got.Equal(tt.want) {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
		})
	}
}
