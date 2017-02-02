package geo

// go test -v github.com/skypies/geo

import "testing"

func TestCentroid(t *testing.T) {
	poly := NewPolygon()
	poly.AddPoint(Latlong{ 0, 0})
	poly.AddPoint(Latlong{ 0,10})
	poly.AddPoint(Latlong{10,10})
	poly.AddPoint(Latlong{10, 0})

	if !poly.Centroid().Equal(Latlong{5,5}) {
		t.Errorf("Centroid not 5,5: %s\n", poly)
	}
}

