package geo

// MapRenderer lets things output themselves as primitive shapes, to be rendered on a map
type MapRenderer interface {
	String() string	
	ToLines() []LatlongLine
	ToCircles() []LatlongCircle
}

// GeoRestricter means things that apply geo restrictions to tracks
type GeoRestrictor interface {
	String() string	

	IntersectsAltitude(int64) bool
	IntersectsLine (LatlongLine) bool
	IntersectsLineDeb (LatlongLine) (bool,string)
	LookForExit() bool  // Should we enter, and then exit, this restrition ? Or is it a window ?

	//IntersectsContiguousLineSet ([]LatlongLine) bool  // O(log N)

	// Also implements MapRenderer
	ToLines() []LatlongLine
	ToCircles() []LatlongCircle
}
