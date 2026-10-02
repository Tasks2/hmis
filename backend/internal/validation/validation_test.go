package validation

import "testing"

func TestRequired(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid value",
			value: "Amy",
			want:  true,
		},
		{
			name:  "empty value",
			value: "",
			want:  false,
		},
		{
			name:  "whitespace",
			value: "   ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Required(tt.value); got != tt.want {
				t.Fatalf("Required(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid email",
			value: "amy@example.com",
			want:  true,
		},
		{
			name:  "missing at sign",
			value: "amyexample.com",
			want:  false,
		},
		{
			name:  "missing domain",
			value: "amy@",
			want:  false,
		},
	}
	//no 63
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Email(tt.value); got != tt.want {
				t.Fatalf("Email(%q) = %v, want %v", tt.name, tt.value, tt.want)
			}
		})
	}
}

func TestDate(t *testing.T) {
	if _, err := Date("2026-10-02"); err != nil {
		t.Fatalf("expected valid date, got error: %v", err)
	}

	if _, err := Date("02-10-2026"); err == nil {
		t.Fatal("expected invalid date")
	}
}

func TestTime(t *testing.T) {
	if _, err := Time("09:30"); err != nil {
		t.Fatalf("expected valid time, got error: %v", err)
	}

	if _, err := Time("9:30 AM"); err == nil {
		t.Fatal("expected invalid time")
	}
}

func TestIs30MinuteSlot(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "on the hour",
			value: "09:00",
			want:  true,
		},
		{
			name:  "half hour",
			value: "09:30",
			want:  true,
		},
		{
			name:  "15 minutes",
			value: "09:15",
			want:  false,
		},
		{
			name:  "invalid time",
			value: "hello",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is30MinuteSlot(tt.value); got != tt.want {
				t.Fatalf(
					"Is30MinuteSlot(%q) = %v, want %v",
					tt.value,
					got,
					tt.want,
				)
			}
		})
	}
}
