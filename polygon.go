package geo

import(
	"fmt"
	pmgeo "github.com/paulmach/go.geo"  // https://godoc.org/github.com/paulmach/go.geo
)

type Polygon struct {
	*pmgeo.Path
	closedPath *pmgeo.Path // transient
}
func NewPolygon() *Polygon { return &Polygon{ Path: pmgeo.NewPath() } }

func (poly *Polygon)String() string {
	return fmt.Sprintf("Poly n=%d, center=%s, avg radius=%.2fKM",
		len(poly.Path.Points()), poly.Centroid(), poly.ApproxRadiusKM())
}

func (poly *Polygon)setClosedPath() {
	if poly.closedPath != nil {
		return
	}
	pts := poly.Path.Points()
	pts = append(pts, pts[0]) // Close the path

	poly.closedPath = pmgeo.NewPath()
	poly.closedPath.PointSet = pts
}

func (poly *Polygon)GetPoints() []Latlong {
	ret := []Latlong{}
	for _,pt := range poly.Path.Points() {
		ret = append(ret, LatlongFromPt(&pt))
	}
	return ret
}

// A slice of sides, in clockwise order, where each side is represented
// as a slice of two points {Start, End}
func (poly *Polygon)GetSides() [][]Latlong {
	poly.setClosedPath()
	pts := poly.closedPath.Points()
	ret := [][]Latlong{}
	for i:=1; i<len(pts); i++ {
		ret = append(ret, []Latlong{LatlongFromPt(&pts[i-1]), LatlongFromPt(&pts[i])})
	}
	return ret
}

func (poly *Polygon)Centroid() Latlong {
	lats,longs := 0.0,0.0

	if n := len(poly.Path.Points()); n > 0 {
		for _,p := range poly.Path.Points() {
			lats += p.Lat()
			longs += p.Lng()
		}
		lats /= float64(n)
		longs /= float64(n)
	}

	return Latlong{lats, longs}
}

// Average distance from centroid, in KM.
func (poly *Polygon)ApproxRadiusKM() float64 {
	c := poly.Centroid()
	dists := 0.0
	for _,pt := range poly.Path.Points() {
		dists += c.DistKM(LatlongFromPt(&pt))
	}
	return dists / float64(len(poly.Path.Points()))
}


// Order matters.
func (poly *Polygon)AddPoint(ll Latlong) {
	poly.Path.PointSet = append(poly.Path.PointSet, *(ll.Pt()))
}

// Note; when a line intersects a vertex, it may be found to intersect lines on both sides,
// and so may contribute that vertex point more than once. So we dedupe.
func (poly *Polygon)IntersectsLine(l LatlongLine) ([]Latlong, bool) {
	poly.setClosedPath()

	ret := []Latlong{}
	pts,_ := poly.closedPath.IntersectionLine(l.Ln())

	deduped := []*pmgeo.Point{}
	for i:=0; i<len(pts); i++ {
		isDupe := false
		for j:=0; j<len(deduped); j++ {
			if pts[i].Equals(deduped[j]) {
				isDupe = true
			}
		}
		if !isDupe { deduped = append(deduped, pts[i]) }
	}

	for _,pt := range deduped {
		// infinite points show up when l is collinear with a segment in the path
		//if math.IsInf(pt.Lat(), 0) || math.IsInf(pt.Lng(), 0) { continue }
		ret = append(ret, LatlongFromPt(pt))
	}
	return ret, (len(ret)>0)
}

// This is *so* similar to LatlongBox.OverlapsLine ...
func (poly *Polygon)OverlapsLine(l LatlongLine) OverlapOutcome {
	// Trivial bounding box test; discard if the line (as a box) has no overlap
	if !poly.Path.Bound().Intersects(l.Bound()) { return Disjoint }
	
	// If either endpoint is in the box, we're containing or straddling.
	// TODO: make this less horrifyingly expensive.
	sInside,eInside := poly.Contains(l.From), poly.Contains(l.To)

	// r2 is the line. If any of it is inside, figure out the line's relation to the box
	if sInside && eInside { return OverlapR2IsContained }
	if sInside            { return OverlapR2StraddlesEnd }
	if eInside            { return OverlapR2StraddlesStart }

	// Points are outside. So, if we intersect, then the line contains the poly.
	if _,intersects := poly.IntersectsLine(l); intersects {
		return OverlapR2Contains
	} else {
		return Disjoint
	}
}

func (p *Polygon)Contains(ll Latlong) bool {
	if !p.Path.Bound().Contains(ll.Pt()) { return false }

  // Contains: build 'ray' from p to (0,0), intersect with path, and
  // expect an odd number of intersections.

	// This is crufty; but the ray end needs to be outside of the polygon.
	// Also, the ray can't be coincident with any side of the polygon; TODO.
	rayEnd := Latlong{0,0}
	for {
		if !p.Path.Bound().Contains(rayEnd.Pt()) { break }
		rayEnd.Lat  += 10
		rayEnd.Long += 11.7
	}

	pts,_ := p.IntersectsLine(ll.LineTo(rayEnd))
	
	n := len(pts)
	rem := n - (n/2)*2 // Crufty int mod(2) remainder
	return (rem == 1)  // if odd, we intersect
}


// Implement MapRenderer
func (poly *Polygon)ToCircles() []LatlongCircle { return nil }
func (poly *Polygon)ToLines() []LatlongLine {
	ret := []LatlongLine{}
	for _,pair := range poly.GetSides() {
		ret = append(ret, pair[0].LineTo(pair[1]))
	}
	return ret
}
