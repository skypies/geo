package geo
// go test -v github.com/skypies/geo

import(
	"testing"
)

func TestNewLatlong(t *testing.T) {
	tests := []struct{
		Input    string
		Expected Latlong
	}{
		{"36.7415306 -121.8942333",         Latlong{36.7415306, -121.8942333}},
		{"36.7415306, -121.8942333",        Latlong{36.7415306, -121.8942333}},
		{"36.7415306 / -121.8942333",       Latlong{36.7415306, -121.8942333}},

		{"36°57\"02.96'N, 121°57\"09.62'W", Latlong{36.9508222, -121.9526722}},

		{"265702.96N,    1115709.62W",     Latlong{26.9508222, -111.9526722}},
		{"465702.96S,    1315709.62E",      Latlong{-46.95082222, 131.95267222}},
	}

	for i,test := range tests {
		actual := NewLatlong(test.Input)
		if !actual.Equal(test.Expected) {
			t.Errorf("[test % 2d] NewLatlong %q: expected %v, got %v", i, test.Input, test.Expected, actual)
		}
	}
}
