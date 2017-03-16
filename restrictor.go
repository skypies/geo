package geo

// MapRenderer lets things output themselves as primitive shapes, to be rendered on a map
type MapRenderer interface {
	String() string
	ToLines() []LatlongLine
	ToCircles() []LatlongCircle
}

// Debugger lets things build debug logs
type Debugger interface {
	Debugf(format string, args ...interface{})
	GetDebug() string // Will delete it.
}

type Intersector interface {
	MapRenderer
	Debugger

	BoundingBox() LatlongBox // 2D bounding box for intersection
	CanContain() bool // Will be false for window intersections
	Contains(Latlong) bool
	OverlapsLine(LatlongLine) OverlapOutcome
	OverlapsAltitude(int64) OverlapOutcome // DisjointR2ComesAfter == val was above the GR
}	

type Restrictor interface {
	Intersector
	IsExclusion() bool // I.e. the restriction means "do not intersect this thing"
	IsNil() bool
}
