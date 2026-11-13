package main

import "fmt"

const (
	substation       = "SUBSTATION"
	transformer      = "TRANSFORMER"
	lvBoard          = "LV_BOARD"
	switchboard      = "SWITCHBOARD"
	switchboardPanel = "SWITCHBOARD_PANEL"

	panelsPerSwitchboard = 5
)

// columns are the columns of an asset file, in order.
var columns = []string{
	"asset_id", "parent_asset_id", "asset_type", "asset_name",
	"operational_status", "voltage_kv", "rating_kva", "manufacturer",
	"model", "serial_number", "commissioned_date",
}

// The values the sample data uses.
var (
	manufacturers = []string{"Arcwell Industries", "Helios Gridworks", "Northstar Electrical", "Meridian Power Systems"}
	ratings       = []string{"1000", "1500", "2000"}
)

// row is one asset of the file being generated.
type row struct {
	id, parentID, assetType, name, status string
	voltage, rating, manufacturer, model  string
	serial, date                          string

	line     int  // line of the CSV file the row is written on
	children int  // how many rows have it as their parent
	used     bool // an error has been put in it
}

func (r *row) record() []string {
	return []string{
		r.id, r.parentID, r.assetType, r.name, r.status, r.voltage,
		r.rating, r.manufacturer, r.model, r.serial, r.date,
	}
}

// file is the asset file being generated.
type file struct {
	hash     string // ends every asset ID, so that files do not clash
	rows     []*row
	findings []finding // what the app should report about the file
}

// buildFile makes a valid asset file of total rows, shaped like the sample
// data: each substation has two transformers, an LV board and switchboards
// of five panels. The assets are shared evenly between the substations.
func buildFile(substations, total int, hash string) *file {
	f := &file{hash: hash}
	start := wordStart(hash)
	others := total - substations

	for i := range substations {
		number := i + 1
		quota := others / substations
		if i < others%substations {
			quota++
		}
		word := substationWord(i, start)

		ss := f.add(nil, &row{
			id:        fmt.Sprintf("SS-%03d", number),
			assetType: substation,
			name:      fmt.Sprintf("%s Substation (%s)", word, regions[(i/2)%len(regions)]),
			voltage:   "22",
			serial:    fmt.Sprintf("SUB-7%03d", number),
			date:      "15/2/20",
		})

		// add puts one more asset under parent, until the quota is used up.
		add := func(parent *row, r *row) *row {
			if quota == 0 {
				return nil
			}
			quota--
			return f.add(parent, r)
		}

		for t := 1; t <= 2; t++ {
			rating := ratings[(i+t)%len(ratings)]
			add(ss, &row{
				id:        fmt.Sprintf("TX-%03d-%d", number, t),
				assetType: transformer,
				name:      fmt.Sprintf("%s Transformer T%d", word, t),
				voltage:   "22",
				rating:    rating,
				model:     "TXR-" + rating,
				serial:    fmt.Sprintf("TX%03d%02dA", number, t),
				date:      "12/3/21",
			})
		}
		add(ss, &row{
			id:        fmt.Sprintf("LVB-%03d-1", number),
			assetType: lvBoard,
			name:      fmt.Sprintf("%s Main LV Board 1", word),
			voltage:   "0.4",
			model:     "LVB-400",
			serial:    fmt.Sprintf("LV%03d01B", number),
			date:      "20/6/21",
		})
		for s := 1; quota > 0; s++ {
			series := seriesLetter(s - 1)
			board := add(ss, &row{
				id:        fmt.Sprintf("SWB-%03d-%d", number, s),
				assetType: switchboard,
				name:      fmt.Sprintf("%s 22kV Switchboard %s", word, series),
				voltage:   "22",
				model:     "SWB-22",
				serial:    fmt.Sprintf("SW%03d%02dC", number, s),
				date:      "10/9/21",
			})
			for p := 1; p <= panelsPerSwitchboard; p++ {
				add(board, &row{
					id:        fmt.Sprintf("PNL-%03d-%d-%02d", number, s, p),
					assetType: switchboardPanel,
					name:      fmt.Sprintf("%s Panel %s%d", word, series, p),
					voltage:   "22",
					model:     "PNL-22",
					serial:    fmt.Sprintf("PN%03d%02d%02d", number, s, p),
					date:      "11/9/21",
				})
			}
		}
	}
	return f
}

// add appends a row under its parent and fills in what every row shares.
func (f *file) add(parent *row, r *row) *row {
	n := len(f.rows)
	r.line = n + 2 // line 1 is the header
	r.id += "-GEN-ID-" + f.hash
	r.serial += "-GEN-SER"

	r.status = "IN_SERVICE"
	if parent != nil {
		r.parentID = parent.id
		parent.children++

		// Only substations are always in service and have no manufacturer.
		r.manufacturer = manufacturers[n%len(manufacturers)]
		switch n % 7 {
		case 3:
			r.status = "MAINTENANCE"
		case 5:
			r.status = "OUT_OF_SERVICE"
		}
	}

	f.rows = append(f.rows, r)
	return r
}

// wordStart picks which names a file starts with from its hash.
func wordStart(hash string) int {
	sum := 0
	for _, c := range hash {
		sum += int(c)
	}
	return sum % wordsPerLetter
}
