package geo

// Various shims to/from the lib we're moving to

import pmgeo "github.com/paulmach/go.geo"

func (ll Latlong)Pt() *pmgeo.Point { return pmgeo.NewPointFromLatLng(ll.Lat,ll.Long) }	
func (l LatlongLine)Ln() *pmgeo.Line { return pmgeo.NewLine(l.From.Pt(), l.To.Pt()) }
func (l LatlongLine)Bound() *pmgeo.Bound { return pmgeo.NewBoundFromPoints(l.From.Pt(), l.To.Pt()) }

func LatlongFromPt(p *pmgeo.Point) Latlong { return Latlong{Lat:p.Lat(), Long:p.Lng()} }
func LatlongLineFromLn(ln pmgeo.Line) LatlongLine {
	return LatlongFromPt(ln.A()).LineTo(LatlongFromPt(ln.B()))
}
func LatlongBoxFromBound(b *pmgeo.Bound) LatlongBox {
	return LatlongFromPt(b.SouthWest()).BoxTo(LatlongFromPt(b.NorthEast()))
}
