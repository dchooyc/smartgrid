package store

import "testing"

func TestEscapeLike(t *testing.T) {
	tests := []struct{ term, want string }{
		{"SS-001", "SS-001"},
		{"100%", `100\%`},
		{"LV_BOARD", `LV\_BOARD`},
		{`a\b`, `a\\b`},
	}
	for _, test := range tests {
		if got := escapeLike(test.term); got != test.want {
			t.Errorf("escapeLike(%q) = %q, want %q", test.term, got, test.want)
		}
	}
}

func TestFormatForDatabase(t *testing.T) {
	if got := formatDate("15/2/09"); got != "2009-02-15" {
		t.Errorf(`formatDate("15/2/09") = %q, want "2009-02-15"`, got)
	}
	if got := formatDate(""); got != "" {
		t.Errorf(`formatDate("") = %q, want ""`, got)
	}
	if got := formatNumber("1e3"); got != "1000" {
		t.Errorf(`formatNumber("1e3") = %q, want "1000"`, got)
	}
	if got := formatNumber("0.4"); got != "0.4" {
		t.Errorf(`formatNumber("0.4") = %q, want "0.4"`, got)
	}
	if got := formatNumber(""); got != "" {
		t.Errorf(`formatNumber("") = %q, want ""`, got)
	}
}
