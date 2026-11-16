package main

import "testing"

func TestWithoutPassword(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"the password is hidden", "postgres://smartgrid:secret@db:5432/smartgrid?sslmode=disable", "postgres://smartgrid:xxxxx@db:5432/smartgrid?sslmode=disable"},
		{"a URL without a password is unchanged", "postgres://smartgrid@db:5432/smartgrid", "postgres://smartgrid@db:5432/smartgrid"},
		{"an unreadable URL is not shown at all", "postgres://smartgrid:sec ret@%zz", "a URL that could not be read"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := withoutPassword(test.url); got != test.want {
				t.Errorf("withoutPassword(%q) = %q, want %q", test.url, got, test.want)
			}
		})
	}
}
