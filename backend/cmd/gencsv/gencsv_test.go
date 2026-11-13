package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"smartgrid/backend/internal/asset"
	"smartgrid/backend/internal/csvfile"
)

// reported runs a generated file through the app's own validation and
// returns the labels it reports for each line, errors and warnings apart.
func reported(t *testing.T, content string) (errs, warnings map[int][]string, rows int) {
	t.Helper()
	_, assets, err := csvfile.Read(strings.NewReader(content))
	if err != nil {
		t.Fatalf("the app rejected the file: %v", err)
	}

	errs, warnings = map[int][]string{}, map[int][]string{}
	for _, e := range asset.ValidateWithStored(assets, nil) {
		errs[e.Line] = append(errs[e.Line], e.Label)
	}
	for _, w := range asset.FindWarnings(assets) {
		warnings[w.Line] = append(warnings[w.Line], w.Label)
	}
	return errs, warnings, len(assets)
}

// expected returns the labels the file's answer sheet promises for each line.
func expected(f *file) (errs, warnings map[int][]string) {
	errs, warnings = map[int][]string{}, map[int][]string{}
	for _, found := range f.findings {
		if found.warning {
			warnings[found.line] = append(warnings[found.line], found.labels...)
		} else {
			errs[found.line] = append(errs[found.line], found.labels...)
		}
	}
	return errs, warnings
}

func describe(byLine map[int][]string) string {
	var lines []int
	for line := range byLine {
		lines = append(lines, line)
	}
	slices.Sort(lines)

	text := ""
	for _, line := range lines {
		labels := slices.Clone(byLine[line])
		slices.Sort(labels)
		text += fmt.Sprintf("\n  line %d: %s", line, strings.Join(labels, "; "))
	}
	if text == "" {
		return " (none)"
	}
	return text
}

// checkAgainstApp generates a file and fails unless the app reports exactly
// what its answer sheet says.
func checkAgainstApp(t *testing.T, o options) {
	t.Helper()
	f := buildFile(o.substations, o.assets, o.hash)
	for i, k := range kinds {
		if o.counts[i] > 0 {
			if err := k.apply(f, o.counts[i]); err != nil {
				t.Fatal(err)
			}
		}
	}
	records, _ := f.table(columnChanges{})
	content, err := writeCSV(records)
	if err != nil {
		t.Fatal(err)
	}

	gotErrs, gotWarnings, rows := reported(t, content)
	wantErrs, wantWarnings := expected(f)
	if rows != o.assets {
		t.Errorf("the file has %d rows, want %d", rows, o.assets)
	}
	if got, want := describe(gotErrs), describe(wantErrs); got != want {
		t.Errorf("errors\n the app reports:%s\n the answer sheet says:%s", got, want)
	}
	if got, want := describe(gotWarnings), describe(wantWarnings); got != want {
		t.Errorf("warnings\n the app reports:%s\n the answer sheet says:%s", got, want)
	}
}

func optionsFor(substations, assets int) options {
	return options{substations: substations, assets: assets, hash: "a1b2c3d4", counts: make([]int, len(kinds))}
}

func TestValidFiles(t *testing.T) {
	sizes := []struct{ substations, assets int }{
		{1, 1}, {1, 16}, {3, 3}, {3, 48}, {3, 50}, {10, 102}, {2, 200}, {60, 700},
	}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("SS_%d_A_%d", size.substations, size.assets), func(t *testing.T) {
			checkAgainstApp(t, optionsFor(size.substations, size.assets))
		})
	}
}

func TestEachKindOfError(t *testing.T) {
	for i, k := range kinds {
		t.Run(k.flag, func(t *testing.T) {
			o := optionsFor(4, 64)
			o.counts[i] = 3

			f := buildFile(o.substations, o.assets, o.hash)
			if err := k.apply(f, 3); err != nil {
				t.Fatal(err)
			}
			if len(f.findings) == 0 {
				t.Fatal("nothing was changed")
			}
			checkAgainstApp(t, o)
		})
	}
}

func TestAllKindsTogether(t *testing.T) {
	o := optionsFor(6, 200)
	for i := range o.counts {
		o.counts[i] = 2
	}
	checkAgainstApp(t, o)
}

func TestTooFewRowsIsAnError(t *testing.T) {
	f := buildFile(1, 2, "a1b2c3d4") // a substation and one transformer
	for _, k := range kinds {
		if err := k.apply(f, 5); err == nil {
			t.Errorf("%s: asking for more rows than there are should fail", k.flag)
		}
	}
}

func TestColumnChanges(t *testing.T) {
	cases := []columnChanges{
		{missing: 1}, {missing: 3}, {duplicate: 2}, {extra: 2},
		{missing: 2, duplicate: 1, extra: 1},
	}
	for _, changes := range cases {
		t.Run(fmt.Sprintf("%+v", changes), func(t *testing.T) {
			records, rejection := buildFile(2, 20, "a1b2c3d4").table(changes)
			content, err := writeCSV(records)
			if err != nil {
				t.Fatal(err)
			}

			_, _, err = csvfile.Read(strings.NewReader(content))
			if err == nil {
				t.Fatal("the app accepted the file")
			}
			if err.Error() != rejection {
				t.Errorf("rejection\n the app says:          %s\n the answer sheet says: %s", err, rejection)
			}
		})
	}
}

func TestNamesAndIDs(t *testing.T) {
	f := buildFile(120, 120, "a1b2c3d4")

	names := map[string]bool{}
	for _, r := range f.rows {
		if !strings.HasSuffix(r.id, "-GEN-ID-a1b2c3d4") {
			t.Errorf("asset ID %q does not end with -GEN-ID-a1b2c3d4", r.id)
		}
		if !strings.HasSuffix(r.serial, "-GEN-SER") {
			t.Errorf("serial number %q does not end with -GEN-SER", r.serial)
		}
		if names[r.name] {
			t.Errorf("the name %q is used twice", r.name)
		}
		names[r.name] = true
	}
	if first := f.rows[0].name; !strings.HasPrefix(first, "M") {
		t.Errorf("the first substation is %q, want a name starting with M", first)
	}
}

func TestLocation(t *testing.T) {
	o := optionsFor(10, 102)
	if folder, _ := o.location(timeAt(1731400000)); folder != "valid" {
		t.Errorf("folder = %q, want valid", folder)
	}

	o.counts[slices.IndexFunc(kinds, func(k kind) bool { return k.flag == "negative-rating" })] = 3
	folder, name := o.location(timeAt(1731400000))
	if folder != "negative-rating" || name != "SS_10_A_102_NEGRATING_3_1731400000_a1b2c3d4" {
		t.Errorf("got %s/%s", folder, name)
	}

	o.columns.missing = 1
	folder, name = o.location(timeAt(1731400000))
	if folder != "mixed" || name != "SS_10_A_102_MCOL_1_NEGRATING_3_1731400000_a1b2c3d4" {
		t.Errorf("got %s/%s", folder, name)
	}
}

func timeAt(unix int64) time.Time { return time.Unix(unix, 0) }
