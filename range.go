package geo

import "time"

type RangeInterface interface {
	Start() float64
	End() float64
}

type OverlapOutcome int
const(
	Undefined OverlapOutcome = iota
	// These two apply to linear (range) overlaps (including altitude), where we have an ordering
	DisjointR2ComesAfter
	DisjointR2ComesBefore
	OverlapR2StraddlesStart
	OverlapR2StraddlesEnd

	// These two apply to all kinds of overlaps
	OverlapR2IsContained
	OverlapR2Contains

	// This two apply to area overlaps, which don't have an ordering
	Disjoint
	OverlapStraddles	

)
func (o OverlapOutcome)String() string {
	return map[OverlapOutcome]string{
		Undefined: "Undefined",
		DisjointR2ComesAfter: "DisjointR2ComesAfter",
		DisjointR2ComesBefore: "DisjointR2ComesBefore",
		OverlapR2StraddlesStart: "OverlapR2StraddlesStart",
		OverlapR2StraddlesEnd: "OverlapR2StraddlesEnd",
		OverlapR2IsContained: "OverlapR2IsContained",
		OverlapR2Contains: "OverlapR2Contains",
		Disjoint: "Disjoint",
		OverlapStraddles: "OverlapStraddles",
	}[o]
}

func (o OverlapOutcome)IsDisjoint() bool {
	return o==DisjointR2ComesBefore || o==DisjointR2ComesAfter || o==Disjoint
}
func (o OverlapOutcome)IsContained() bool {
	return o==OverlapR2IsContained || o==OverlapR2Contains
}

// When they're identical, prefer r1 to contain r2
func RangeOverlap(r1, r2 RangeInterface) OverlapOutcome {
	if r1.Start() > r2.End() {
		return DisjointR2ComesBefore
	} else if r1.End() < r2.Start() {
		return DisjointR2ComesAfter
	} else if r2.End() > r1.End() && r2.Start() > r1.Start() {
		return OverlapR2StraddlesEnd
	} else if r2.End() < r1.End() && r2.Start() < r1.Start() {
		return OverlapR2StraddlesStart
	} else if r2.Start() >= r1.Start() {
		return OverlapR2IsContained
	} else {
		return OverlapR2Contains
	}
}
	
type Float64Range struct {
	U,V float64
}
func (r Float64Range)Start() float64 { return r.U }
func (r Float64Range)End() float64 { return r.V }

type Int64Range struct {
	U,V int64
}
func (r Int64Range)Start() float64 { return float64(r.U) }
func (r Int64Range)End() float64 { return float64(r.V) }

type TimeRange struct {
	U,V time.Time
}
func time2float64(t time.Time) float64 {
	return float64(t.Unix()) + (float64(t.UnixNano()) / 1000000000.0)
}
func (r TimeRange)Start() float64 { return time2float64(r.U) }
func (r TimeRange)End() float64 { return time2float64(r.V) }


