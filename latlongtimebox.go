package geo

import (
	"fmt"
	"time"
)

// It is a latlong box, with times; used to decide if tracks overlap
type LatlongTimeBox struct {
	LatlongBox
	Start,End time.Time
}

func (tbox LatlongTimeBox)String() string {
	return fmt.Sprintf("%s-%s, %s[+%s]", tbox.SW, tbox.NE,
		tbox.Start.Format("15:04:05.999999"), tbox.End.Sub(tbox.Start))
}

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

func (tb1 LatlongTimeBox)SpaceCompare(tb2 LatlongTimeBox) bool {
	return true
}

// }}}
// {{{ tbox.CompareBoxeSlices

// Given two contiguous sequences of boxes, do they overlap in time
// and space well enough to be the same thing ?

// NOTE: should precede this with a boundingbox test; tracks that can plausibly glue together
// but which don't actually overlap in time will return 'false' from this.

// overlaps: if we should consider them the same thing
// conf: how confident we are
// debug: some debug text about it.
func CompareBoxSlices(b1,b2 *[]LatlongTimeBox) (overlaps bool, conf float64, debug string) {
	// Counters:
	//  for a box: how many from other track had time-overlap; how many also had space-overlap
	//  for a track: how many boxes had time-space-overlap; how many had just time-overlap
	
	
	// track A is the one with the earliest timestamp.
	tA,tB := b1,b2  // Need to do an init flip

	debug = fmt.Sprintf("* tA: %s\n* tB: %s\n", AsTimeline(*tA, (*tA)[0].Start),
		AsTimeline(*tB, (*tA)[0].Start))
	
	// fast-forward track B until it overlaps.
	iA,iB := 0,0
	
	// Then walk up both, zip-zagging and comparing ...
	nFlips := 0
	for {
		if iA>=len(*tA) || iB>=len(*tB) {
			// All done. (account for uncompared boxes ?)
			break
		}

		// We have just flipped (or initialized) to looking at track A.
		debug += fmt.Sprintf("----------------------------- iter [%d,%d]\n", iA, iB)
		
/*
		// Any B-boxes we encompass: compare them (increment boxscores for A&B), and finalize them
		while currA.FullyEnclose(nextB) {
			timeOverlap,spaceOverlap := currA.Compare(nextB)
			// Add partial scores to currA
			// Fully score nextB; finalize it
			nextB++
		}

*/

		if false {
			// boxes line up exactly - increment both sides, don't flip
			
			// Here: it is possible that both tracks break at exactly the same
			// point. If so, we should not compare currA against nextB. Just
			// increment A, and do not flip.
			//if nextB.start >= currA.end {
				//currA++

		} else {
			// nextB is not fully encompassed; thus it is the last box currA touches
			//timeOverlap,spaceOverlap := tA[iA].Compare(tB[iB])
			space,cDebug := LoggedCompare(nFlips, tA, tB, iA, iB, "--- ")
			debug += cDebug
			_=space

			// Add scores to currA; finalize it
			// Add prelim scores to nextB (which becomes currA :>)

			// FLIP !!
			iA,iB = iB,iA+1
			tA,tB = tB,tA
			nFlips++
			
			//tA,tB = tB,tA
			//newCurrA = nextB // Preserve the prelim scores above
			//newNextB = currA + 1
			//currA,nextB = newCurrA,newNextB
		}
	}

	return
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
