package geo

import "fmt"

// Window is a finite line across the land, with altitude restrictions
// It implements GeoRestricter and MapRenderer.
type Window struct {
	LatlongLine  // embedded
	MinAltitude float64
	MaxAltitude float64
}

func (w Window)String() string {
	return fmt.Sprintf("GeoRestrict{Window}: %s-%s %.1fKM, [%.0f,%.0f]ft",
		w.From, w.To, w.To.DistKM(w.From), w.MinAltitude, w.MaxAltitude)
}

// Implement GeoRestricter interface
func (w Window)LookForExit() bool { return false }
func (w Window)IntersectsLine(l LatlongLine) bool {
	_,intersects := w.LatlongLine.Intersects(l)
	return intersects
}
func (w Window)IntersectsAltitude(alt int64) bool {
	if w.MinAltitude > 0 && float64(alt) < w.MinAltitude { return false }
	if w.MaxAltitude > 0 && float64(alt) > w.MaxAltitude { return false }
	return true
}

func (w Window)IntersectsLineDeb(l LatlongLine) (bool,string) { return w.IntersectsLine(l),"" }

// Implement MapRenderer interface
func (w Window)ToCircles() []LatlongCircle { return []LatlongCircle{} }
func (w Window)ToLines()   []LatlongLine   { return []LatlongLine{w.LatlongLine} }
