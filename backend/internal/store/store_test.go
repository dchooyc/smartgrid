package store

import "testing"

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
