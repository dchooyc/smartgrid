// gencsv generates sample asset CSV files for testing the app by hand: valid
// ones, and ones broken in a chosen way. Each file comes with an answer sheet
// saying what the app should report for it.
//
// Run it from the backend folder, e.g.
//
//	go run ./cmd/gencsv -substations 3 -assets 48 -negative-rating 2
//
// or answer the questions asked by scripts/generate-csv.sh.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// options are what the command-line arguments ask for.
type options struct {
	substations int
	assets      int // rows in the file, substations included
	columns     columnChanges
	counts      []int  // rows to break, for each of kinds
	hash        string // ends every asset ID
	out         string // folder the files are sorted into
}

func main() {
	var o options
	list := flag.Bool("list", false, "list the kinds of row error and exit")
	flag.IntVar(&o.substations, "substations", 3, "number of substations")
	flag.IntVar(&o.assets, "assets", 0, "number of rows, substations included (default 16 per substation)")
	flag.IntVar(&o.columns.missing, "missing-columns", 0, "columns to remove")
	flag.IntVar(&o.columns.duplicate, "duplicate-columns", 0, "columns to repeat")
	flag.IntVar(&o.columns.extra, "extra-columns", 0, "unknown columns to add")
	flag.StringVar(&o.hash, "hash", "", "8 characters that end every asset ID (default random)")
	flag.StringVar(&o.out, "out", filepath.Join("..", "generated-csv"), "folder to put the files in")
	o.counts = make([]int, len(kinds))
	for i, k := range kinds {
		flag.IntVar(&o.counts[i], k.flag, 0, "rows to break: "+strings.ToLower(k.title))
	}
	flag.Parse()

	if *list {
		for _, k := range kinds {
			fmt.Printf("%s\t%s\n", k.flag, k.title)
		}
		return
	}

	csvPath, sheetPath, err := run(o, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "gencsv:", err)
		os.Exit(1)
	}
	fmt.Println("CSV file:     ", csvPath)
	fmt.Println("Answer sheet: ", sheetPath)
}

// run generates the file and its answer sheet and returns where they are.
func run(o options, now time.Time) (csvPath, sheetPath string, err error) {
	if o.assets == 0 {
		o.assets = 16 * o.substations
	}
	if o.hash == "" {
		random := make([]byte, 4)
		if _, err := rand.Read(random); err != nil {
			return "", "", err
		}
		o.hash = hex.EncodeToString(random)
	}
	if err := o.check(); err != nil {
		return "", "", err
	}

	folder, name := o.location(now)
	content, sheet, err := generate(o, name+".csv")
	if err != nil {
		return "", "", err
	}

	directory := filepath.Join(o.out, folder)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", "", err
	}
	csvPath = filepath.Join(directory, name+".csv")
	sheetPath = filepath.Join(directory, name+"_answers.txt")
	if err := os.WriteFile(csvPath, []byte(content), 0o644); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(sheetPath, []byte(sheet), 0o644); err != nil {
		return "", "", err
	}
	return csvPath, sheetPath, nil
}

// generate returns the text of the CSV file and of its answer sheet.
func generate(o options, csvName string) (content, sheet string, err error) {
	f := buildFile(o.substations, o.assets, o.hash)
	for i, k := range kinds {
		if o.counts[i] == 0 {
			continue
		}
		if err := k.apply(f, o.counts[i]); err != nil {
			return "", "", err
		}
	}

	records, rejection := f.table(o.columns)
	if content, err = writeCSV(records); err != nil {
		return "", "", err
	}
	return content, f.answerSheet(csvName, o.substations, rejection), nil
}

func (o options) check() error {
	switch {
	case o.substations < 1:
		return errors.New("there must be at least 1 substation")
	case o.assets < o.substations:
		return fmt.Errorf("%d assets cannot include %d substations: ask for at least as many assets as substations", o.assets, o.substations)
	case o.columns.missing < 0 || o.columns.duplicate < 0 || o.columns.extra < 0:
		return errors.New("the numbers of columns cannot be negative")
	case o.columns.missing >= len(columns):
		return fmt.Errorf("at most %d columns can be removed", len(columns)-1)
	case o.columns.duplicate > len(columns)-o.columns.missing:
		return fmt.Errorf("only %d columns are left to repeat", len(columns)-o.columns.missing)
	case len(o.hash) != 8:
		return errors.New("the hash must be 8 characters")
	}
	for i, count := range o.counts {
		if count < 0 {
			return fmt.Errorf("-%s cannot be negative", kinds[i].flag)
		}
	}
	return nil
}

// location returns the folder a file is sorted into, named after the one
// thing wrong with it, and the file's name without an extension, e.g.
// "negative-rating" and "SS_3_A_48_NEGRATING_2_1731400000_a1b2c3d4".
func (o options) location(now time.Time) (folder, name string) {
	var folders, codes []string
	add := func(count int, folder, code string) {
		if count > 0 {
			folders = append(folders, folder)
			codes = append(codes, fmt.Sprintf("%s_%d", code, count))
		}
	}
	add(o.columns.missing, "missing-columns", "MCOL")
	add(o.columns.duplicate, "duplicate-columns", "DCOL")
	add(o.columns.extra, "extra-columns", "XCOL")
	for i, k := range kinds {
		add(o.counts[i], k.flag, k.code)
	}

	switch len(folders) {
	case 0:
		folder = "valid"
	case 1:
		folder = folders[0]
	default:
		folder = "mixed"
	}
	parts := []string{fmt.Sprintf("SS_%d_A_%d", o.substations, o.assets)}
	parts = append(parts, codes...)
	parts = append(parts, fmt.Sprint(now.Unix()), o.hash)
	return folder, strings.Join(parts, "_")
}
