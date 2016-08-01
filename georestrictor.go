package geo

//// https://play.golang.org/p/jF-SdC6Lxf

// MapRenderer lets things output themselves as primitive shapes, to be rendered on a map
type MapRenderer interface {
	String() string
	ToLines() []LatlongLine
	ToCircles() []LatlongCircle
}

// An Intersector thing implements standard intersection routines.
type Intersector interface {
	MapRenderer
	IntersectsAltitude(int64) bool
	IntersectsLine (LatlongLine) bool
	IntersectsLineDeb (LatlongLine) (bool,string)
	//IntersectsContiguousLineSet ([]LatlongLine) bool  // O(log N), plz
	LookForExit() bool  // Should we enter, and then exit, this restrition ? Or is it a window ?
}

// A Restrictor thing implements intersection, and logic
type Restrictor interface {
	Intersector
	MustNotIntersect() bool
}

// Wrapper objects, to turn geo objects into restrictors
type LatlongBoxRestrictor struct {
	LatlongBox
	MustNotIntersectVal bool
}
func (r LatlongBoxRestrictor)MustNotIntersect() bool { return r.MustNotIntersectVal }

type LatlongCircleRestrictor struct {
	LatlongCircle
	MustNotIntersectVal bool
}
func (r LatlongCircleRestrictor)MustNotIntersect() bool { return r.MustNotIntersectVal }

type WindowRestrictor struct {
	Window
	MustNotIntersectVal bool
}
func (r WindowRestrictor)MustNotIntersect() bool { return r.MustNotIntersectVal }
