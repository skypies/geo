package geo

import (
	"fmt"
	"math"
	"time"
)

// It is a latlong box, with times; used to decide if tracks overlap
type LatlongTimeBox struct {
	LatlongBox   // embedded
	Start,End    time.Time
	HeadingDelta float64 // How much the heading altered during this box
	I,J          int     // range of indices of trackpoints that generated the box

	// These values are used to decide if a box is too approximate to be safe to compare
	Source               string // too hard to backtrack when debugging, just say so here
	Interpolated         bool // The box was interpolated, so might be pretty bogus
	RunLength            int  // how many boxes were in this run of interpolation
	CentroidHeadingDelta float64  // This is fiddly; see below
	Debug                string
}

// {{{ tbox.String

func (tbox LatlongTimeBox)String() string {
	str := fmt.Sprintf("{%3d,%3d} %s+%5.4f,%5.4f, %-12.12s[+%s], %3.0fdeg [%s]", tbox.I, tbox.J,
		tbox.SW, (tbox.NE.Lat - tbox.SW.Lat), (tbox.NE.Long - tbox.SW.Long), 
		tbox.Start.Format("15:04:05.999"), tbox.End.Sub(tbox.Start), tbox.HeadingDelta, tbox.Source)
	if tbox.Interpolated {
		str += fmt.Sprintf(" InterpDelta:%3.0fdeg, n=%d", tbox.CentroidHeadingDelta, tbox.RunLength)
	}
	return str
}

// }}}

// {{{ AsTimeline

func AsTimeline(boxes []LatlongTimeBox, start time.Time) string {
	return AsFlaggedTimeline(boxes,start,-1)
}

// }}}
// {{{ AsFlaggedTimeline

func AsFlaggedTimeline(boxes []LatlongTimeBox, start time.Time, flagIndex int) string {
	if len(boxes) == 0 { return "[empty]" }

	str := ""
	prefix := boxes[0].Start.Sub(start).Seconds()
	for prefix>0 { str += " "; prefix-- }
	str += "|"

	prevEnd := boxes[0].Start
	for i,box := range boxes {
		gap := box.Start.Sub(prevEnd).Seconds()
		len := box.End.Sub(box.Start).Seconds()
		for gap>0 { str += " "; gap-- }

		char := "-"
		if i == flagIndex { char = "=" }
		for len>1 { str += char; len-- }
		str += "|"
		prevEnd = box.End
	}
	return str
}

// }}}
// {{{ LoggedCompare

func LoggedCompare(nFlips int, tA,tB *[]LatlongTimeBox, iA,iB int, prefix string) (space bool, deb string) {

	if true {
		// This logging stuff is exponentially expensive, don't use on real flights
		return (*tA)[iA].SpaceCompare((*tB)[iB]), ""
	}
	
	if nFlips%2 == 0 {
		tStart := (*tA)[0].Start
		deb += fmt.Sprintf("%s o-tA: %s\n", prefix, AsFlaggedTimeline(*tA, tStart, iA))
		deb += fmt.Sprintf("%s o-tB: %s\n", prefix, AsFlaggedTimeline(*tB, tStart, iB))
	} else {
		tStart := (*tB)[0].Start
		deb += fmt.Sprintf("%s x-tB: %s\n", prefix, AsFlaggedTimeline(*tB, tStart, iB))
		deb += fmt.Sprintf("%s x-tA: %s\n", prefix, AsFlaggedTimeline(*tA, tStart, iA))
	}	

	space = (*tA)[iA].SpaceCompare((*tB)[iB])
	// deb += fmt.Sprintf("%s space=%v\n", prefix, space)
	
	return
}

// }}}

// {{{ tbox.Enclose

func (tb *LatlongTimeBox)Enclose(p Latlong, t time.Time) {
	tb.LatlongBox.Enclose(p)
	if tb.Start.After(t) {
		tb.Start = t
	}
	if tb.End.Before(t) {
		tb.End = t
	}
}

// }}}
// {{{ tbox.EnsureMinSide

// If the box is too narrow or too tall, fatten it out
func (tbox *LatlongTimeBox)EnsureMinSide(min float64) {
	if tbox.LongWidth() < min {
		c := tbox.Center()
		tbox.SW.Long = c.Long - min/2.0
		tbox.NE.Long = c.Long + min/2.0
	}
	if tbox.LatHeight() < min {
		c := tbox.Center()
		tbox.SW.Lat = c.Lat - min/2.0
		tbox.NE.Lat = c.Lat + min/2.0
	}
}

// }}}
// {{{ tbox.SpaceCompare

// Do the boxes overlap in latlong space ?
func (tb1 LatlongTimeBox)SpaceCompare(tb2 LatlongTimeBox) bool {
	disp,amount := tb1.LatlongBox.OverlapsWith(tb2.LatlongBox)
	_=amount // Ignore

	return ! disp.IsDisjoint()
}

// }}}
// {{{ tbox.AltitudeSpeedHeadingSimilar

// Do the boxes have similar altitudes, speed and heading ?
const (
	maxAltitudeDeviation    = 100.0
	//maxHeadingDeviation     = 2.0
	maxGroundspeedDeviation = 4.0
)
func (tb1 LatlongTimeBox)AltitudeSpeedHeadingSimilar(tb2 LatlongTimeBox) bool {
	return false

	// We need trackpoints for this !!
	
	//if math.Abs(b.CentroidHeadingDelta) > maxCentroidDeviation { return true }

}

// }}}
// {{{ tbox.TooApproximateForComparison

// 1. If a box implies a large change in heading (HeadingDeviation) - or
// the centre of the box is off to the side of where the heading
// currently points (CentroidDeviation) - then the interpolation is
// likely ropey, so we can't usefully compare the box for a time/space
// overlap.
//
// 2. If the box is part of an interpolation run of >2 boxes, then
// bail altogether. It turned out that tracks with different amounts
// of interpolation could have pretty linear (non-curvey)
// interpolation tracks, but which didn't entirely overlap thanks to
// the 'sawtooth' nature of the boxes leaving some smaller boxes
// orphaned. So when comparing against a run of boxes, give up.
const (
	maxCentroidDeviation = 40.0
	maxHeadingDeviation = 40.0
)
func (b LatlongTimeBox)TooApproximateForComparison() bool {
	if b.Interpolated {
		if math.Abs(b.CentroidHeadingDelta) > maxCentroidDeviation { return true }
		if math.Abs(b.HeadingDelta) > maxHeadingDeviation { return true }
		if b.RunLength > 2 { return true } // harsh ... but fair
	}
	return false
}

// }}}

// {{{ BoxResult{}, BoxSliceResult{}

type BoxResult struct {
	NumTimeOverlaps      int   // How many boxes, from the other track, overlapped in time
	NumTimeSpaceOverlaps int   // How many boxes, from the other track, overlapped in time and space
	NumIgnores           int   // How many comparisons we ignored, because of sketchy interpolation
	TimeOverlapIndices []int   // Which boxes from the other track we have Time-only overlaps
}

type BoxSliceResult struct {
	B                        []BoxResult
	NumBoxes                   int
	NumBoxesJustTimeOverlap    int
	NumBoxesTimeSpaceOverlap   int
	NumBoxesIgnored            int
  BadBoxIndices            []int   // For debugging; which boxes had just time overlap
	BadBoxOtherIndices     [][]int   // And for each such box, which boxes in other track it messed
}

func (r BoxSliceResult)String() string {
	str := fmt.Sprintf("--{N=%d, t=%d, t+s=%d, i=%d}--\n", r.NumBoxes,
		r.NumBoxesJustTimeOverlap, r.NumBoxesTimeSpaceOverlap, r.NumBoxesIgnored)
	for i,br := range r.B {
		str += fmt.Sprintf(" [%3d] t=%2d t+s=%2d\n", i,
			br.NumTimeSpaceOverlaps, br.NumTimeSpaceOverlaps)
	}
	return str
}

func NewResults(b *[]LatlongTimeBox) BoxSliceResult {
	return BoxSliceResult{
		B: make([]BoxResult, len(*b)),
		NumBoxes: len(*b),
		BadBoxIndices: []int{},
		BadBoxOtherIndices: [][]int{},
	}
}

func (r *BoxSliceResult)Finalize() {
	r.NumBoxesJustTimeOverlap,r.NumBoxesTimeSpaceOverlap = 0,0
	for i,br := range r.B {
		if br.NumTimeSpaceOverlaps>0 {
			r.NumBoxesTimeSpaceOverlap++
		} else if br.NumTimeOverlaps>0 {
			r.NumBoxesJustTimeOverlap++
			r.BadBoxIndices = append(r.BadBoxIndices, i)
			r.BadBoxOtherIndices = append(r.BadBoxOtherIndices, br.TimeOverlapIndices)
		}
		if br.NumIgnores > 0 {
			r.NumBoxesIgnored++
		}
	}
}

// }}}
// {{{ compareAndScore

func compareAndScore(b1,b2 LatlongTimeBox, i1,i2 int, r1,r2 *BoxResult) string {
	str := ""

	if b1.TooApproximateForComparison() || b2.TooApproximateForComparison() {
		r1.NumIgnores++
		r2.NumIgnores++
		return str
	}


	if b1.SpaceCompare(b2) {
		r1.NumTimeSpaceOverlaps++
		r2.NumTimeSpaceOverlaps++
	} else {
		// No space overlap.
		// If the altitudes and headings are close, and the space-gap isn't too bad, write it off
		// as clock skew (and abandon) ... TBD, need trackpoints to get the data !
		// if b1.AltitudeSpeedHeadingSimilar(b2) {...}

		// str += fmt.Sprintf(" * * * SPONPONPON [%d:%d] * *\n", i1, i2)
		
		r1.TimeOverlapIndices = append(r1.TimeOverlapIndices, i2)
		r2.TimeOverlapIndices = append(r2.TimeOverlapIndices, i1)
	}

	// We're here because the boxes have a time overlap. Without a corresponding space overlap,
	// this will kill the whole compare, so wait until we haven't bailed on the comparison before
	// recording this.
	r1.NumTimeOverlaps++
	r2.NumTimeOverlaps++
	
	return str
}

// }}}
// {{{ CompareBoxeSlicesNew

// Given two contiguous sequences of boxes (ordered in time), do they
// overlap in time and space well enough to be considered parts of the
// same thing ?

// Basic idea: find where boxes overlap in time, and then see if they
// also overlap in space. If there are lots of overlaps in both space
// and time, it's a good match. If there are any boxes that overlap in
// time, but *don't* overlap in space, then they are divergent.

// We perform a single pass other both slices, jumping from one to the
// other, so that all boxes get compared with any time-overlapping
// boxes from the other slice.

// overlaps: if we should consider them the same thing
// conf: how confident we are
// str: some debug text about it.
func CompareBoxSlicesNew(b1,b2 *[]LatlongTimeBox) (overlaps bool, conf float64, str string) {
	if len(*b1)==0 || len(*b2) == 0 { return false, 0.0, "at least one slice was empty, abort" }
	
	r1 := NewResults(b1)
	r2 := NewResults(b2)
	rA,rB := r1,r2
	tA,tB := b1,b2

	// track A will be the one with the earliest timestamp
	if (*b2)[0].Start.Before((*b1)[0].Start) {
		tA,tB = tB,tA
		rA,rB = rB,rA
	}

	str = fmt.Sprintf("* t1: %s --> %s\n* t2: %s --> %s\n",
		(*b1)[0].Start, (*b1)[len(*b1)-1].End, (*b2)[0].Start, (*b2)[len(*b2)-1].End)

	// If no overlap at all, just bail
	if (*tA)[len(*tA)-1].End.Before((*tB)[0].Start) {
		return false, 0.0, str + "* no time overlap at all, bailing\n"
	}
	
	iA,iB := 0,0
	// Track A starts first; fast-forward it until there is an overlap in time.
	for (*tA)[iA].End.Before((*tB)[iB].Start) {
		iA++
	}
	
	// Now start the main loop
	debug := false
	nFlips := 0

outerLoop:
	for {
		// Every time we start at the top of the loop, the tracks and indices are set up like this:
		//  tA:  . . . |=====| . . .   (i.e. tA[iA], "currA", is a box
		//  tB:   . . . . |==. . . .   (i.e. tB[iB], "currB", is also a box, that starts after currA
		//                              but we don't know where it ends
		// or, maybe:
		//  tA:  . . . |=====| . . . 
		//  tB:   . . .|==. . . .      (i.e. both boxes start at the exact same point.

		//if iA == 20 || iB == 20 { debug = true } else { debug = false }
		
		if iA>=len(*tA) || iB>=len(*tB) {
			str += fmt.Sprintf("** Breaking loop; iA=%d, iB=%d\n", iA, iB)
			// All done. No need to account for uncompared boxes; the results structs will zero out.
			break outerLoop
		}
		
		if debug {str += fmt.Sprintf("----------------------------- iter [%d,%d]\n", iA, iB)}

		// Does the current B box end before the current A box ? I.e. do we fully enclose it ?
		for (*tB)[iB].End.Before( (*tA)[iA].End ) {
			// tA:    |========|-----|-----|-----|
			// tB:       |=|--
			if(debug){str += "* enclosed case; curr-A fully encloses curr-B\n"}
			str += compareAndScore((*tA)[iA], (*tB)[iB], iA, iB, &(rA.B[iA]), &(rB.B[iB]))

			// Move along to the next box in B
			iB++
			if iB>=len(*tB) { break outerLoop }
			if(debug){str += fmt.Sprintf("-.-.-.-.-.-.-.-.-.-.-.-.-.-.- iter [%d,%d]\n", iA, iB)}
		}
		
		// So, we don't (or no longer) enclose curr-B; it ends after curr-A ends.
		
		// If curr-B *starts* at exactly the same place that curr-A ends, then no
		// need to compare (they don't overlap) or to flip; just move A along one.
		if (*tB)[iB].Start.Equal( (*tA)[iA].End ) {
			// tA:    |========|-----|-----|-----|
			// tB:             |===|
			if(debug){str += "* degenerate case (perfect alignment) - \n"}
			iA++
			if iA>=len(*tA) { break outerLoop}

		} else {
			// tA:    |========|-----|-----|-----|
			// tB:         |======|
			if debug {str += "* default case; A and B have partial overlap, B straddles end of A\n"}
			str += compareAndScore((*tA)[iA], (*tB)[iB], iA, iB, &(rA.B[iA]), &(rB.B[iB]))

			// Move along one box on track A ...
			iA++
			if iA>=len(*tA) { break outerLoop}

			// ... and FLIP !
			if(debug){str += "* * * * * * * * * * FLIP * * * * * * * * * *   "}
			iA,iB = iB,iA
			tA,tB = tB,tA
			rA,rB = rB,rA
			nFlips++
		}
	}

	// OK, we're done zig-zagging and comparing. Let's see what each slice has to say for itself.
	r1.Finalize()
	r2.Finalize()

	str += "**** Outcome\n"

	if r1.NumBoxesJustTimeOverlap > 0 || r2.NumBoxesJustTimeOverlap > 0 {
		// Boxes that overlap in time, but not in space, are bad. For now, tolerate none of these.
		overlaps, conf = false, 0.0

		for j,boxIndex := range r1.BadBoxIndices {
			str += fmt.Sprintf("* r1 bad box %3d - %s\n%s", boxIndex, (*b1)[boxIndex],
				(*b1)[boxIndex].Debug )
			for _,otherBoxIndex := range r1.BadBoxOtherIndices[j] {
				str += fmt.Sprintf("**  otherbox %3d - %s\n%s", otherBoxIndex, (*b2)[otherBoxIndex],
					(*b2)[otherBoxIndex].Debug )
			}
		}
		for j,boxIndex := range r2.BadBoxIndices {
			str += fmt.Sprintf("* r2 bad box %3d - %s\n%s", boxIndex, (*b2)[boxIndex],
				(*b2)[boxIndex].Debug )
			for _,otherBoxIndex := range r2.BadBoxOtherIndices[j] {
				str += fmt.Sprintf("**  otherbox %3d - %s\n%s", otherBoxIndex, (*b1)[otherBoxIndex],
					(*b1)[otherBoxIndex].Debug )
			}
		}

		str += "* some time-only overlaps found, rejecting\n"
	} else {
		if r1.NumBoxesTimeSpaceOverlap > 0 || r2.NumBoxesTimeSpaceOverlap > 0 {
			str += "* time+space overlaps!\n"
			overlaps = true
			r1Overlap := float64(r1.NumBoxesTimeSpaceOverlap) / float64(r1.NumBoxes)
			r2Overlap := float64(r2.NumBoxesTimeSpaceOverlap) / float64(r2.NumBoxes)
			conf = (r1Overlap + r2Overlap) / 2.0
		} else {
			str += "* no overlap in time at all (bounding box test, plz)\n"
		}
	}

	//str += fmt.Sprintf("* overlaps=%v (confidence=%.2f)\n", overlaps, conf)
	//str += fmt.Sprintf("***** Results\n* Track 1: %s* Track 2: %s", r1, r2)
	
	str += fmt.Sprintf("tA[%d:+%d-%d(%d?)],tB[%d:+%d-%d(%d?)]=%.2f,%v",
		r1.NumBoxes, r1.NumBoxesTimeSpaceOverlap, r1.NumBoxesJustTimeOverlap, r1.NumBoxesIgnored,
		r2.NumBoxes, r2.NumBoxesTimeSpaceOverlap, r2.NumBoxesJustTimeOverlap, r2.NumBoxesIgnored,
		conf, overlaps)
		
	return
}

// }}}

// {{{ CompareBoxeSlicesOld

func boxesShouldBeCompared(b1, b2 LatlongTimeBox) bool {
	return !(b1.TooApproximateForComparison() || b2.TooApproximateForComparison())
}

// Given two contiguous sequences of boxes (ordered in time), do they
// overlap in time and space well enough to be considered parts of the
// same thing ?

// Basic idea: find where boxes overlap in time, and then see if they
// also overlap in space. If there are lots of overlaps in both space
// and time, it's a good match. If there are any boxes that overlap in
// time, but *don't* overlap in space, then they are divergent.

// We perform a single pass other both slices, jumping from one to the
// other, so that all boxes get compared with any time-overlapping
// boxes from the other slice.

// overlaps: if we should consider them the same thing
// conf: how confident we are
// str: some debug text about it.
func CompareBoxSlicesOld(b1,b2 *[]LatlongTimeBox) (overlaps bool, conf float64, str string) {
	if len(*b1)==0 || len(*b2) == 0 { return false, 0.0, "at least one slice was empty, abort" }
	
	r1 := NewResults(b1)
	r2 := NewResults(b2)
	rA,rB := r1,r2
	tA,tB := b1,b2

	// track A will be the one with the earliest timestamp
	if (*b2)[0].Start.Before((*b1)[0].Start) {
		tA,tB = tB,tA
		rA,rB = rB,rA
	}

	str = fmt.Sprintf("* t1: %s --> %s\n* t2: %s --> %s\n",
		(*b1)[0].Start, (*b1)[len(*b1)-1].End, (*b2)[0].Start, (*b2)[len(*b2)-1].End)
	//	str = fmt.Sprintf("* tA: %s\n* tB: %s\n", AsTimeline(*tA, (*tA)[0].Start),
  //		AsTimeline(*tB, (*tA)[0].Start))

	// If no overlap at all, just bail
	if (*tA)[len(*tA)-1].End.Before((*tB)[0].Start) {
		return false, 0.0, str + "* no time overlap at all, bailing\n"
	}
	
	iA,iB := 0,0
	// Track A starts first; fast-forward it until there is an overlap in time.
	for (*tA)[iA].End.Before((*tB)[iB].Start){
		iA++
	}
	
	// Now start the main loop
	debug := false
	nFlips := 0

outerLoop:
	for {
		// Every time we start at the top of the loop, the tracks and indices are set up like this:
		//  tA:  . . . |=====| . . .   (i.e. tA[iA], "currA", is a box
		//  tB:   . . . . |==. . . .   (i.e. tB[iB], "currB", is also a box, that starts after currA
		//                              but we don't know where it ends
		// or, maybe:
		//  tA:  . . . |=====| . . . 
		//  tB:   . . .|==. . . .      (i.e. both boxes start at the exact same point.

		//if iA == 20 || iB == 20 { debug = true } else { debug = false }
		
		if iA>=len(*tA) || iB>=len(*tB) {
			str += fmt.Sprintf("** Breaking loop; iA=%d, iB=%d\n", iA, iB)
			// All done. No need to account for uncompared boxes; the results structs will zero out.
			break outerLoop
		}
		
		if debug {str += fmt.Sprintf("----------------------------- iter [%d,%d]\n", iA, iB)}

		// Does the current B box end before the current A box ? I.e. do we fully enclose it ?
		for (*tB)[iB].End.Before( (*tA)[iA].End ) {
			// tA:    |========|-----|-----|-----|
			// tB:       |=|--
			if(debug){str += "* enclosed case; curr-A fully encloses curr-B\n"}

			// When we're dealing with tracks that have interpolated boxes along curves - and perhaps
			// taken crazy shortcuts - we should just ignore the comparison.
			scoreThis := boxesShouldBeCompared((*tA)[iA], (*tB)[iB])
			if scoreThis {
				rA.B[iA].NumTimeOverlaps++
				rB.B[iB].NumTimeOverlaps++
			
				// Now space-compare iA with iB (log scores in both).
				spaceOverlap,_ := LoggedCompare(nFlips, tA, tB, iA, iB, "--- ")
				if spaceOverlap {
					rA.B[iA].NumTimeSpaceOverlaps++
					rB.B[iB].NumTimeSpaceOverlaps++
				} else {
					// No space overlap. Track the failed comparison for later.
					rA.B[iA].TimeOverlapIndices = append(rA.B[iA].TimeOverlapIndices, iB)
					rB.B[iB].TimeOverlapIndices = append(rB.B[iB].TimeOverlapIndices, iA)
					//str += "================={ DOOM2 }======================\n"
					//str += fmt.Sprintf("= currA: %d, %s\n= currB: %d, %s\n", iA, (*tA)[iA], iB, (*tB)[iB]
				}
			} else {
				if(debug){str += fmt.Sprintf("* skipped a compare, [%d/%d]\n", iA, iB)}
				rA.B[iA].NumIgnores++
				rB.B[iB].NumIgnores++
			}

			// Move along to the next box in B
			iB++
			if iB>=len(*tB) { break outerLoop }
			if(debug){str += fmt.Sprintf("-.-.-.-.-.-.-.-.-.-.-.-.-.-.- iter [%d,%d]\n", iA, iB)}
		}
		
		// So, we don't (or no longer) enclose curr-B; it ends after curr-A ends.
		
		// If curr-B *starts* at exactly the same place that curr-A ends, then no
		// need to compare (they don't overlap) or to flip; just move A along one.
		if (*tB)[iB].Start.Equal( (*tA)[iA].End ) {
			// tA:    |========|-----|-----|-----|
			// tB:             |===|
			if(debug){str += "* degenerate case (perfect alignment) - \n"}
			iA++
			if iA>=len(*tA) { break outerLoop}

		} else {
			// tA:    |========|-----|-----|-----|
			// tB:         |======|
			if debug {
				str += "* default case; A and B have partial overlap, and B straddles the end of A\n"
			}

			scoreThis := boxesShouldBeCompared((*tA)[iA], (*tB)[iB])
			if scoreThis {			
				rA.B[iA].NumTimeOverlaps++
				rB.B[iB].NumTimeOverlaps++

				// Now space-compare iA with iB (log scores in both).
				spaceOverlap,_ := LoggedCompare(nFlips, tA, tB, iA, iB, "--- ")
				if spaceOverlap && scoreThis {
					rA.B[iA].NumTimeSpaceOverlaps++
					rB.B[iB].NumTimeSpaceOverlaps++
				} else {
					// No space overlap. Track the failed comparison for later.
					rA.B[iA].TimeOverlapIndices = append(rA.B[iA].TimeOverlapIndices, iB)
					rB.B[iB].TimeOverlapIndices = append(rB.B[iB].TimeOverlapIndices, iA)
					//str += "================={ DOOM }======================\n"
					//str += fmt.Sprintf("= currA: %d, %s\n= currB: %d, %s\n", iA, (*tA)[iA], iB, (*tB)[iB])
				}
			} else {
				rA.B[iA].NumIgnores++
				rB.B[iB].NumIgnores++
				if(debug){str += fmt.Sprintf("* skipped a compare, [%d/%d]\n", iA, iB)}
			}

			// Move along one box on track A ...
			iA++
			if iA>=len(*tA) { break outerLoop}

			// ... and FLIP !
			if(debug){str += "* * * * * * * * * * FLIP * * * * * * * * * *   "}
			iA,iB = iB,iA
			tA,tB = tB,tA
			rA,rB = rB,rA
			nFlips++
		}
	}

	// OK, we're done zig-zagging and comparing. Let's see what each slice has to say for itself.
	r1.Finalize()
	r2.Finalize()

	str += "**** Outcome\n"

	if r1.NumBoxesJustTimeOverlap > 0 || r2.NumBoxesJustTimeOverlap > 0 {
		// Boxes that overlap in time, but not in space, are bad. For now, tolerate none of these.
		overlaps, conf = false, 0.0

		for j,boxIndex := range r1.BadBoxIndices {
			str += fmt.Sprintf("* r1 bad box %3d - %s\n%s", boxIndex, (*b1)[boxIndex],
				(*b1)[boxIndex].Debug )
			for _,otherBoxIndex := range r1.BadBoxOtherIndices[j] {
				str += fmt.Sprintf("**  otherbox %3d - %s\n%s", otherBoxIndex, (*b2)[otherBoxIndex],
					(*b2)[otherBoxIndex].Debug )
			}
		}
		for j,boxIndex := range r2.BadBoxIndices {
			str += fmt.Sprintf("* r2 bad box %3d - %s\n%s", boxIndex, (*b2)[boxIndex],
				(*b2)[boxIndex].Debug )
			for _,otherBoxIndex := range r2.BadBoxOtherIndices[j] {
				str += fmt.Sprintf("**  otherbox %3d - %s\n%s", otherBoxIndex, (*b1)[otherBoxIndex],
					(*b1)[otherBoxIndex].Debug )
			}
		}

		str += "* some time-only overlaps found, rejecting\n"
	} else {
		if r1.NumBoxesTimeSpaceOverlap > 0 || r2.NumBoxesTimeSpaceOverlap > 0 {
			str += "* time+space overlaps!\n"
			overlaps = true
			r1Overlap := float64(r1.NumBoxesTimeSpaceOverlap) / float64(r1.NumBoxes)
			r2Overlap := float64(r2.NumBoxesTimeSpaceOverlap) / float64(r2.NumBoxes)
			conf = (r1Overlap + r2Overlap) / 2.0
		} else {
			str += "* no overlap in time at all (bounding box test, plz)\n"
		}
	}

	//str += fmt.Sprintf("* overlaps=%v (confidence=%.2f)\n", overlaps, conf)
	//str += fmt.Sprintf("***** Results\n* Track 1: %s* Track 2: %s", r1, r2)
	
	str += fmt.Sprintf("tA[%d:+%d-%d(%d?)],tB[%d:+%d-%d(%d?)]=%.2f,%v",
		r1.NumBoxes, r1.NumBoxesTimeSpaceOverlap, r1.NumBoxesJustTimeOverlap, r1.NumBoxesIgnored,
		r2.NumBoxes, r2.NumBoxesTimeSpaceOverlap, r2.NumBoxesJustTimeOverlap, r2.NumBoxesIgnored,
		conf, overlaps)
		
	return
}

// }}}

func CompareBoxSlices(b1,b2 *[]LatlongTimeBox) (overlaps bool, conf float64, str string) {
	return CompareBoxSlicesNew(b1,b2)
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
