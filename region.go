package geo

import()

type Region interface {
	String() string
	
	ContainsPoint (Latlong) bool
	IntersectsBox (LatlongBox) bool
	IntersectsLine (LatlongLine) bool

	// Render the region into primitive shapes, to display on maps etc.
	ToLines() []LatlongLine
	ToCircles() []LatlongCircle
}
