package geo

import(
	"fmt"
)

const (
	kKmPerLatitudeDegreeAtSFO  = 111.2   // (36,-122)->(37,-122) == 111 KM (heading north)
	kKmPerLongitudeDegreeAtSFO = 88.08   // (36,-122)->(36,-121) ==  88 KM (heading east)
)

type LatlongBox struct {
	SW, NE       Latlong
	Floor, Ceil  int64  // altitude, feet; zero means "don't care". Nonzero means >= or <=, depending
}
func (box LatlongBox)String() string {
	str := fmt.Sprintf("%s-%s", box.SW, box.NE)
	if box.Floor > 0 || box.Ceil > 0 {
		str += fmt.Sprintf("[%d,%d]", box.Floor, box.Ceil)
	}
	return str
}

// Derive the other two corners on demand
func (box LatlongBox)SE() Latlong { return Latlong{Lat:box.SW.Lat , Long:box.NE.Long} }
func (box LatlongBox)NW() Latlong { return Latlong{Lat:box.NE.Lat , Long:box.SW.Long} }

// Derive bounded lines for the sides
func (box LatlongBox)BottomSide() LatlongLine { return box.SW.LineTo(box.SE()) }
func (box LatlongBox)LeftSide()   LatlongLine { return box.SW.LineTo(box.NW()) }
func (box LatlongBox)TopSide()    LatlongLine { return box.NW().LineTo(box.NE) }
func (box LatlongBox)RightSide()  LatlongLine { return box.SE().LineTo(box.NE) }

func (box LatlongBox)LongWidth() float64 { return box.NE.Long - box.SW.Long }
func (box LatlongBox)LatHeight() float64 { return box.NE.Lat - box.SW.Lat }

func (box LatlongBox)Center() Latlong {
	return Latlong{
		Lat: (box.SW.Lat + box.NE.Lat) / 2.0,
		Long: (box.SW.Long + box.NE.Long) / 2.0,
	}
}

// Returns a box, centred on ll, that is of size (width,height)
func (ll Latlong)Box(widthKm,heightKm float64) LatlongBox {
	// This is a hack; these constants are only valid quite close to SFO itself.
	latOffset := (heightKm / kKmPerLatitudeDegreeAtSFO) / 2.0
	longOffset := (widthKm / kKmPerLongitudeDegreeAtSFO) / 2.0
	return LatlongBox{
		SW: Latlong{ll.Lat-latOffset, ll.Long-longOffset},
		NE: Latlong{ll.Lat+latOffset, ll.Long+longOffset},
	}
}

func (from Latlong)BoxTo(to Latlong) LatlongBox {
	box := LatlongBox{}
	if (from.Lat < to.Lat) {
		box.SW.Lat, box.NE.Lat = from.Lat, to.Lat
	} else {
		box.SW.Lat, box.NE.Lat = to.Lat, from.Lat
	}
	if (from.Long < to.Long) {
		box.SW.Long, box.NE.Long = from.Long, to.Long
	} else {
		box.SW.Long, box.NE.Long = to.Long, from.Long
	}
	return box
}

func (box LatlongBox)Contains(pos Latlong) bool {
	if (pos.Long < box.SW.Long) { return false }
	if (pos.Lat  < box.SW.Lat ) { return false }
	if (pos.Long > box.NE.Long) { return false }
	if (pos.Lat  > box.NE.Lat ) { return false }
	return true
}

// Enclose increases the sixe of the box to include the point, if it doesn't fit
func (box *LatlongBox)Enclose(pos Latlong) {
	if (pos.Long < box.SW.Long) { box.SW.Long = pos.Long }
	if (pos.Lat  < box.SW.Lat ) { box.SW.Lat = pos.Lat }
	if (pos.Long > box.NE.Long) { box.NE.Long = pos.Long }
	if (pos.Lat  > box.NE.Lat ) { box.NE.Lat = pos.Lat }
}

func (box LatlongBox)LatRange() Float64Range { return Float64Range{box.SW.Lat, box.NE.Lat} }
func (box LatlongBox)LongRange() Float64Range { return Float64Range{box.SW.Long, box.NE.Long} }

// Returned float *should* be the fraction of b1 that overlaps with b2
func (b1 LatlongBox)OverlapsWith(b2 LatlongBox) (OverlapOutcome,float64) {
	latDisp := RangeOverlap(b1.LatRange(), b2.LatRange())
	longDisp := RangeOverlap(b1.LongRange(), b2.LongRange())
	
	if latDisp.IsDisjoint() || longDisp.IsDisjoint() {
		return Disjoint, 0.0
	} else if latDisp == longDisp {
		switch latDisp {
		case OverlapR2IsContained: return OverlapR2IsContained, 1.0
		case OverlapR2Contains:    return OverlapR2Contains, 1.0
		}
	}

	return OverlapStraddles, 1.0
}

// Implement the Restrictor interface
func (box LatlongBox)LookForExit() bool { return true }
/*
func (box LatlongBox)IntersectsLine(l LatlongLine) bool {
	// Trivial bounding box test; discard if the line (as a box) has no overlap
	if !box.IntersectsBox(l.Box()) { return false }

	// If the box contains either point, we're good.
	if box.Contains(l.From) || box.Contains(l.To) { return true }

	// Else: we know the boxes overlap, but both line points are outside of it; so ensure
	// that the line has a (bounded) intersection with the box.
	if _,isect := box.BottomSide().Intersects(l); isect { return true }
	if _,isect := box.LeftSide().Intersects(l); isect { return true }
	if _,isect := box.RightSide().Intersects(l); isect { return true }
	if _,isect := box.TopSide().Intersects(l); isect { return true }
	
	return false
}
*/
func (box LatlongBox)IntersectsLine(l LatlongLine) bool {
	return ! box.OverlapsLine(l).IsDisjoint()
}

func (box LatlongBox)OverlapsLine(l LatlongLine) OverlapOutcome {
	// Trivial bounding box test; discard if the line (as a box) has no overlap
	if !box.IntersectsBox(l.Box()) { return Disjoint }

	// If either endpoint is in the box, we're containing or straddling.
	sInside,eInside := box.Contains(l.From), box.Contains(l.To)

	// r2 is the line. If any of it is inside, figure out the line's relation to the box
	if sInside && eInside { return OverlapR2IsContained }
	if sInside            { return OverlapR2StraddlesEnd }
	if eInside            { return OverlapR2StraddlesStart }
	
	// Else: we know the boxes overlap, but both line points are outside of it; if the line
	// has a (bounded) intersection with any edge of the box, then we deem the box to be
	// contained by the line.
	if _,isect := box.BottomSide().Intersects(l); isect { return OverlapR2Contains }
	if _,isect := box.LeftSide().Intersects(l); isect { return OverlapR2Contains }
	if _,isect := box.RightSide().Intersects(l); isect { return OverlapR2Contains }
	if _,isect := box.TopSide().Intersects(l); isect { return OverlapR2Contains }
	
	return Disjoint
}


func (box LatlongBox)IntersectsAltitude(alt int64) bool {
	if box.Floor > 0 && alt < box.Floor { return false }
	if box.Ceil > 0  && alt > box.Ceil  { return false }
	return true
}

func (box LatlongBox)IntersectsLineDeb(l LatlongLine) (bool, string) {
	return box.IntersectsLine(l),""
}

// Implement Region interface (defunct ?)
func (box LatlongBox)ContainsPoint(pos Latlong) bool {
	return box.Contains(pos)
}
func (b1 LatlongBox)IntersectsBox(b2 LatlongBox) bool {
	outcome,_ := b1.OverlapsWith(b2)
	return ! outcome.IsDisjoint()
}

// Implement MapRenderer interface
func (b LatlongBox)ToLines() []LatlongLine {
	return []LatlongLine{ b.BottomSide(), b.LeftSide(), b.RightSide(), b.TopSide() }
}
func (c LatlongBox)ToCircles() []LatlongCircle { return []LatlongCircle{} }


// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
