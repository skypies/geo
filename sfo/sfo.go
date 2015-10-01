package sfo

import "github.com/skypies/geo"

// A bunch of constants relating to SFO

var (
	KLatlongSFO = geo.Latlong{37.6188172, -122.3754281}
	KLatlongSJC = geo.Latlong{37.3639472, -121.9289375}	
	KLatlongSERFR1 = geo.Latlong{37.221516, -121.992987} // This is the centerpoint for maps viewport

	KBoxSFO120K = KLatlongSFO.Box(80,80)  // This is the box in which we look for new flights
	
	KFixes = map[string]geo.Latlong{
		"SERFR": geo.Latlong{36.0683056, -121.3646639},
		"NRRLI": geo.Latlong{36.4956000, -121.6994000},
		"WWAVS": geo.Latlong{36.7415306, -121.8942333},
		"EPICK": geo.Latlong{36.9508222, -121.9526722},
		"EDDYY": geo.Latlong{37.3264500, -122.0997083},
		"SWELS": geo.Latlong{37.3681556, -122.1160806},
		"MENLO": geo.Latlong{37.4636861, -122.1536583},
		"WPOUT": geo.Latlong{37.1194861, -122.2927417},
		"THEEZ": geo.Latlong{37.5034694, -122.4247528},
		"WESLA": geo.Latlong{37.6643722, -122.4802917},
		"MVRKK": geo.Latlong{37.7369722, -122.4544500},
	}

	SFOClassBMap = geo.ClassBMap{
		Name: "SFO",
		Center: KLatlongSFO,
		Sectors: []geo.ClassBSector{
			// Magnetic declination at SFO: 13.68
			geo.ClassBSector{
				StartBearing: 0,
				EndBearing: 360,
				Steps: []geo.Cylinder{
					{ 7,  0, 100},   // from origin to  7NM : 100/00 (no floor)
					{10, 15, 100},   // from   7NM  to 10NM : 100/15
					{15, 30, 100},   // from  10NM  to 15NM : 100/30
					{20, 40, 100},   // from  15NM  to 20NM : 100/40
					{25, 60, 100},   // from  20NM  to 25NM : 100/60
					{30, 80, 100},   // from  25NM  to 30NM : 100/80
				},
			},
			// ... more sectors go here !
		},
	}
)

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
