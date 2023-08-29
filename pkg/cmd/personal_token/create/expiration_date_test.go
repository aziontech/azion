package create

import (
	"testing"
	"time"
)

// TestParseExpirationDate cases "1d", "2w", "2m", "1y", "18/08/2023", "2023-02-12"
func TestParseExpirationDate(t *testing.T) {
	type args struct {
		currentDate      time.Time
		expirationString string
	}

	currentDate := time.Date(2008, 11, 01, 0, 0, 0, 0, &time.Location{})

	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "1 day", args: args{currentDate, "1d"},
			want: time.Date(2008, 11, 02, 0, 0, 0, 0, &time.Location{})},
		{
			name: "2 week", args: args{currentDate, "2w"},
			want: time.Date(2008, 11, 15, 0, 0, 0, 0, &time.Location{}),
		},
		{
			name: "3 months", args: args{currentDate, "3m"},
			want: time.Date(2009, 01, 30, 0, 0, 0, 0, &time.Location{}),
		},
		{
			name: "1 year", args: args{currentDate, "1y"},
			want: time.Date(2009, 11, 01, 0, 0, 0, 0, &time.Location{}),
		},
		{
			name: "Format Br", args: args{currentDate, "04/11/2008"},
			want: time.Date(2008, 11, 04, 0, 0, 0, 0, &time.Location{}),
		},
		{
			name: "Format DB", args: args{currentDate, "2008-11-04"},
			want: time.Date(2008, 11, 04, 0, 0, 0, 0, &time.Location{}),
		},
		{
			name: "Format DB - T", args: args{currentDate, "2008-11-04T00:00"},
			want: time.Date(2008, 11, 04, 0, 0, 0, 0, &time.Location{}),
		},
		{
			name: "Incorrect format", args: args{currentDate, "2008-11"},
			want: time.Time{}, wantErr: true,
		},
		{
			name: "Date shorter than the current one", args: args{currentDate, "2008-10-01"},
			want: time.Time{}, wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseExpirationDate(tt.args.currentDate, tt.args.expirationString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExpirationDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Day() != tt.want.Day() || got.Month() != tt.want.Month() || got.Year() != tt.want.Year() {
				t.Errorf("ParseExpirationDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
