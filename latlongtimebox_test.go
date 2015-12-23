package geo
// go test -v github.com/skypies/geo

import(
	"fmt"
	"testing"
	"time"
)

var (
	spaceBox = LatlongBox{SW:Latlong{0,0}, NE:Latlong{10,10}}
	timeSeqs = [][][]int{
/*
		// Basic flip
		{ {0,    5,   10 },               // A: |----|----|
			{   2,    7    }, },            // B:   |----|
*/
		// Longer flip
		{ {0,    5,   10,   15 },            // A: |----|----|----|
			{   2,   7,     12}, },            // B:   |----|----|

		/*
		// We fully contain
		{ {0,    5,           13    },    // A: |----|--------|
			{   2,     7,   10,     15}, }, // B:   |----|----|----|
		// Tracks perfectly align
		{ {0,    5,   10    },            // A: |-----|----|
			{   2, 5,      12 }, },         // B:   |---|------|
*/
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

func TestCompareBox(t *testing.T) {
	t.Errorf("BAH")
}

func TestCompareBoxSlice(t *testing.T) {
	for _,vals := range timeSeqs {
		b1 := timeSeqToBoxSlice(vals[0])
		b2 := timeSeqToBoxSlice(vals[1])
		overlaps,conf,debug := CompareBoxSlices(&b1,&b2)
		fmt.Printf("** Debug (%v,%.2f):-\n%s", overlaps, conf, debug)
	}
	t.Errorf("BAH")
}


/*
func TestDistAlongLine(t *testing.T) {
	for i,vals := range distalongline {
		lFrom,lTo,pos := Latlong{vals[0],vals[1]}, Latlong{vals[2],vals[3]}, Latlong{vals[4],vals[5]}
		line := lFrom.BuildLine(lTo)
		actual := line.DistAlongLine(pos)
		expected := vals[6]
		if math.Abs(actual-expected) > 0.001 {
			t.Errorf("[%d] distalongline was %f, expected %f", i, actual, expected)
		}
		fmt.Printf("Line:%s, pos:%s, dist:%.3f\n", line, pos, actual)
	}
}
*/

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
