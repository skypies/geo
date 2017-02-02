package geo

import(
	"fmt"
	pmgeo "github.com/paulmach/go.geo"
)

type Polygon struct {
	*pmgeo.Path
}

func NewPolygon() *Polygon {
	return &Polygon{ Path: pmgeo.NewPath() }
}


func (poly *Polygon)String() string {
	return fmt.Sprintf("Poly n=%d, center=%s", len(poly.Path.Points()), poly.Centroid())
}

func (poly *Polygon)Centroid() Latlong {
	lat,long := 0.0,0.0

	if len(poly.Path.Points()) > 0 {
		for _,p := range poly.Path.Points() {
			lat += p.Lat()
			long += p.Lng()
		}
		lat /= float64(len(poly.Path.Points()))
		long /= float64(len(poly.Path.Points()))
	}

	return Latlong{lat, long}
}

func (poly *Polygon)AddPoint(ll Latlong) {
	p := pmgeo.Point{}
	p.SetLat(ll.Lat)
	p.SetLng(ll.Long)
	poly.Path.PointSet = append(poly.Path.PointSet, p)
}


// Implement GeoRestricter interface
// Implement MapRenderer interface

/*

type Polygon struct {
	P []Latlong // Points must be in order, clockwise, and describe a convex polygon
}

func NewPolygon() *Polygon {
	return &Polygon{ P: []Latlong{} }
}

func (poly *Polygon)String() string {
	return fmt.Sprintf("Poly n=%d, center=%s", len(poly.P), poly.Centroid())
}

func (poly *Polygon)Centroid() Latlong {
	lat,long := 0.0,0.0
	for _,p := range poly.P {
		lat += p.Lat
		long += p.Long
	}
	if len(poly.P) > 0 {
		lat /= float64(len(poly.P))
		long /= float64(len(poly.P))
	}
	return Latlong{lat, long}
}

func (poly *Polygon)AddPoint(p Latlong) {
	poly.P = append(poly.P, p)
}

*/
