package geo
// go test -v github.com/skypies/geo

import(
	"testing"
	"time"
)

var (
	spaceBox = LatlongBox{SW:Latlong{0,0}, NE:Latlong{10,10}}
	timeSeqs = [][][]int{
/*
		// Basic flip
		{ {0,    5,   10 },                      // A: |----|----|
			{   2,    7    },                      // B:   |----|
			{1} },
			
		// Longer flip
		{ {0,    5,   10,   15 },                // A: |----|----|----|
			{   2,   7,     12},                   // B:   |----|----|
			{1} },

		// We fully contain
		{ {0,    5,           13    },           // A: |----|--------|
			{   2,     7,   10,     15},           // B:   |----|----|----|
			{1} },

		// We fully contain multiples
		{ {0,    5,                14,   18  },  // A: |----|----------|---|
			{   2,     7,   10,   13,  15},        // B:   |----|---|--|---|
			{1} },

		// Tracks perfectly align (short and long)
		{ {0,    5,   10    },                   // A: |-----|----|
			{   2, 5,      12 },                   // B:   |---|------|
			{1} },
		{ {0,    5,   10,    14    },            // A: |-----|----|---|
			{   2, 5,      12,    16 },            // B:   |---|------|---|
			{1} },

		// No overlap !
		{ {0,    5,   10            },           // A: |-----|----|
			{               12,    16 },           // B:              |---|
			{0},}, // false

		// Some under-run
		{ {0,    5,   10,    14,     18},        // A: |-----|----|---|---|
			{               12,    16 },           // B:             |---|
			{1} }, // true
*/
		// Some over-run
		{ {0,    5,   10,    14,     18},        // A: |-----|----|---|---|
			{   3,    8                  },         // B:    |---|
			{1} }, // true
	}
)

func timeSeqToBoxSlice(k []int) []LatlongTimeBox {
	boxes := []LatlongTimeBox{}
	t,_ := time.Parse("2006.01.02 15:03:04", "2100.01.01 08:00:00")
	for i:=1; i<len(k); i++ {
		box := LatlongTimeBox{
			LatlongBox: spaceBox,
			Start:      t.Add(time.Second * time.Duration(k[i-1])),
			End:        t.Add(time.Second * time.Duration(k[i])),
		}
		boxes = append(boxes, box)
	}
	return boxes
}

func TestCompareBoxSliceBasicZigzag(t *testing.T) {
	for i,vals := range timeSeqs {
		b1 := timeSeqToBoxSlice(vals[0])
		b2 := timeSeqToBoxSlice(vals[1])
		expected := (vals[2][0] > 0)
		overlaps,conf,debug := CompareBoxSlices(&b1,&b2)
		//fmt.Printf("** Debug (%v,%.2f):-\n%s\n", overlaps, conf, debug)
		if overlaps != expected {
			t.Errorf("%,1f\n%s\n[% d] overlap said %v, expected %v", conf, debug, i, overlaps, expected)
		}
	}
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
