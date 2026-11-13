package main

import "fmt"

// The sample data names its substations Aurora to Lighthouse, so generated
// ones carry on from M. Each letter has a few words, so that going round the
// alphabet again gives new names.
var nameWords = [][]string{
	{"Maple", "Marlin", "Meadow", "Mistral"},
	{"Nimbus", "Northgate", "Nectar", "Nova"},
	{"Oakridge", "Opal", "Orchard", "Osprey"},
	{"Pinnacle", "Pebble", "Prairie", "Pioneer"},
	{"Quartz", "Quay", "Quill", "Quarry"},
	{"Ridgeway", "Raven", "Riverbend", "Rowan"},
	{"Summit", "Sapphire", "Sandstone", "Spruce"},
	{"Tidewater", "Topaz", "Timber", "Thistle"},
	{"Umber", "Upland", "Unity", "Union"},
	{"Vista", "Valley", "Verdant", "Vanguard"},
	{"Willow", "Westwind", "Waterfront", "Wren"},
	{"Xenon", "Xylo", "Xanadu", "Xeric"},
	{"Yarrow", "Yellowstone", "Yew", "Yonder"},
	{"Zenith", "Zephyr", "Zircon", "Zinnia"},
}

const wordsPerLetter = 4

var regions = []string{"North", "East", "South", "West"}

// substationWord names the substation at index (from 0). start says which
// word of each letter the first trip round the alphabet uses, so that
// different files get different names.
func substationWord(index, start int) string {
	letters := len(nameWords)
	round := index / letters
	word := nameWords[index%letters][(start+round)%wordsPerLetter]

	// Every word has been used once: tell the repeats apart by number.
	if repeat := round / wordsPerLetter; repeat > 0 {
		return fmt.Sprintf("%s %d", word, repeat+1)
	}
	return word
}

// seriesLetter names the nth (from 0) switchboard of a substation: A, B, ...
func seriesLetter(n int) string {
	if n < 26 {
		return string(rune('A' + n))
	}
	return fmt.Sprintf("S%d", n+1)
}
