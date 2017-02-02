package geo

// go test -v github.com/skypies/geo

import "testing"

type LineTest struct {
	p1,p2,p3,p4 Latlong
	o           bool
}

func TestLineObservers(t *testing.T) {
	// Remember; latlongs are (y,x) not (x,y)
	pOrigin := Latlong{10,10}
	pHoriz  := Latlong{10,20}
	pVert   := Latlong{20,10}
	pDiag   := Latlong{20,20}

	if pOrigin.LineTo(pHoriz).IsVertical() { t.Errorf("horiz triggered vert") }
	if pOrigin.LineTo(pDiag).IsVertical() { t.Errorf("diag triggered vert") }
	if !pOrigin.LineTo(pVert).IsVertical() { t.Errorf("vert didn't trigger vert") }

	if !pOrigin.LineTo(pOrigin).IsDegenerate() {
		t.Errorf("not degenerate: %s", pOrigin.LineTo(pOrigin))
	}
	if pOrigin.LineTo(pDiag).IsDegenerate() { t.Errorf("degenerate") }
}

func TestLineIntersects(t *testing.T) {
	tests := []LineTest{
		{Latlong{1.0,1.0}, Latlong{1.0,4.0}, Latlong{2.0,1.0}, Latlong{2.0,4.0}, false },  // parallel
		{Latlong{1.0,0.0}, Latlong{1.0,3.0}, Latlong{0.0,2.0}, Latlong{3.0,2.0}, true },   // cross
		{Latlong{0.0,0.0}, Latlong{6.0,3.0}, Latlong{0.0,3.0}, Latlong{6.0,0.0}, true },   // diagcross
		{Latlong{1.0,0.0}, Latlong{1.0,3.0}, Latlong{0.0,6.0}, Latlong{3.0,6.0}, false },  // disjoint
		{Latlong{0.0,0.0}, Latlong{5.0,0.0}, Latlong{0.0,0.0}, Latlong{0.0,5.0}, true },   // touch
		{Latlong{0.0,6.0}, Latlong{6.0,6.0}, Latlong{6.0,0.0}, Latlong{6.0,6.0}, true },   // touch
	}

	for i,test := range tests {
		l1, l2 := test.p1.LineTo(test.p2), test.p3.LineTo(test.p4)
		
		actualPos,actualOutcome := l1.Intersects(l2)
		if actualOutcome != test.o {
			t.Errorf("[t%d] exepected %v, got %v {%s} {%s} [%s]", i, test.o, actualOutcome, l1,l2,actualPos)
		}
	}
}

func TestWhichSide(t *testing.T) {
	tests := []struct{
		A,B,C  Latlong
		Out    int
	}{
		// Line.From, Line.To, Point, outcome
		{Latlong{5,0}, Latlong{5,10}, Latlong{ 0,5}, -1},
		{Latlong{5,0}, Latlong{5,10}, Latlong{10,5}, +1},
		{Latlong{5,0}, Latlong{5,10}, Latlong{ 5,5},  0},
		{Latlong{5,10}, Latlong{5,0}, Latlong{ 0,5}, +1},
		{Latlong{5,10}, Latlong{5,0}, Latlong{10,5}, -1},
		{Latlong{5,10}, Latlong{5,0}, Latlong{ 5,5},  0},

		{Latlong{0,0}, Latlong{10,-10}, Latlong{ 0,  0},  0},
		{Latlong{0,0}, Latlong{10,-10}, Latlong{ 5, -5},  0},
		{Latlong{0,0}, Latlong{10,-10}, Latlong{10,-10},  0},
		{Latlong{0,0}, Latlong{10,-10}, Latlong{ 5,  0}, -1},
		{Latlong{0,0}, Latlong{10,-10}, Latlong{ 5,-10}, +1},

		{Latlong{0,0}, Latlong{10,-10}, Latlong{-1, 20}, -1},
	}

	for i,test := range tests {
		line := test.A.LineTo(test.B)
		outcome := line.WhichSide(test.C)

		if test.Out != outcome {
			t.Errorf("[t%d] wanted %d, got %d", i, test.Out, outcome)
		}
	}
}
