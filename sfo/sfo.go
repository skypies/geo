package sfo

import(
	"net/http"
	"sort"
	"github.com/skypies/geo"
)

// A bunch of constants relating to SFO

var (
	// Retire these three.
	KLatlongSFO = geo.Latlong{37.6188172, -122.3754281}
	KLatlongSJC = geo.Latlong{37.3639472, -121.9289375}	
	KLatlongSERFR1 = geo.Latlong{37.221516, -121.992987} // This is the centerpoint for maps viewport

	KBoxSnarfingCatchment = KLatlongSFO.Box(125,125)  // The box in which we look for new flights

	// Boxes used in a few reports
	KBoxSFO10K = KLatlongSFO.Box(12,12)
	KBoxPaloAlto20K = geo.Latlong{37.433536,-122.1310187}.Box(6,7)

	KAirports = map[string]geo.Latlong{
		"KSFO": geo.Latlong{37.6188172, -122.3754281},
		"KSJC": geo.Latlong{37.3639472, -121.9289375},
		"KOAK": geo.Latlong{37.7212597, -122.2211489},
	}
							 
	// http://www.myaviationinfo.com/FixState.php?FixState=CALIFORNIA
	KFixes = map[string]geo.Latlong{
		// SERFR2 & WWAVS1
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

		// BRIXX
		"CORKK": geo.Latlong{37.7335889, -122.4975500},
		"BRIXX": geo.Latlong{37.6178444, -122.3745278},
		"LUYTA": geo.Latlong{37.2948889, -122.2045528},
		"JILNA": geo.Latlong{37.2488056, -122.1495000},
		"YADUT": geo.Latlong{37.2039889, -122.0232778},

		// BIGSURTWO
		"CARME": geo.Latlong{36.4551833, -121.8797139},
		"ANJEE": geo.Latlong{36.7462861, -121.9648917},
		"SKUNK": geo.Latlong{37.0075944, -122.0332278},
		"BOLDR": geo.Latlong{37.1708861, -122.0761667},

		// BDEGA2
		"LOZIT": geo.Latlong{37.899325, -122.673194},
		"BGGLO": geo.Latlong{38.224589, -122.767506},
		"GEEHH": geo.Latlong{38.453333, -122.428650},
		"MSCAT": geo.Latlong{38.566697, -122.671667},
		"JONNE": geo.Latlong{38.551042, -122.863275},
		"AMAKR": geo.Latlong{39.000000, -123.750000},
		"DEEAN": geo.Latlong{38.349164, -123.302289},
		"MRRLO": geo.Latlong{38.897547, -122.578233},
		"MLBEC": geo.Latlong{38.874772, -122.958989},

		// Things for SFO arrivals
		"HEMAN": geo.Latlong{37.5338500, -122.1733333},
		"NEPIC": geo.Latlong{37.5858944, -122.2968833},

		// Things for SFO departures		
		"PORTE": geo.Latlong{37.4897861, -122.4745778},
		"SSTIK": geo.Latlong{37.6783444, -122.3616583},

		// Things for Oceanic
		"PPEGS": geo.Latlong{37.3920722, -122.2817222},
		"ALLBE": geo.Latlong{37.5063889, -127.0000000},
		"ALCOA": geo.Latlong{37.8332528, -125.8345250},
		"CINNY": geo.Latlong{36.1816667, -124.7600000},
		"PAINT": geo.Latlong{38.0000000, -125.5000000},
		"OSI"  : geo.Latlong{37.3925000, -122.2813000},
		"PIRAT": geo.Latlong{37.2576500, -122.8633528},
		"PYE"  : geo.Latlong{38.0797567, -122.8678275},
		"STINS": geo.Latlong{37.8236111, -122.7566667},
		"HADLY": geo.Latlong{37.4022222, -122.5755556},
		
		"PONKE": geo.Latlong{37.4588167, -121.9960528},
		"WETOR": geo.Latlong{37.4847194, -122.0571417},

		// Things for SJC/SILCN3
		"VLLEY": geo.Latlong{36.5091667, -121.4402778},
		"GUUYY": geo.Latlong{36.7394444, -121.5411111},
		"SSEBB": geo.Latlong{36.9788889, -121.6425000},
		"GSTEE": geo.Latlong{37.0708333, -121.6716667},
		"KLIDE": geo.Latlong{37.1641667, -121.7130556},
		"BAXBE": geo.Latlong{36.7730556, -121.6263889},
		"APLLE": geo.Latlong{37.0338889, -121.8050000},

		// Randoms
		"PARIY": geo.Latlong{37.3560056, -121.9231222},  // SJC ?
		"ZORSA": geo.Latlong{37.3627583, -122.0500306},
		
		// Things for East Bay
		"HOPTA": geo.Latlong{37.78501944, -122.154},
		"BOYSS": geo.Latlong{38.02001944, -122.3778639},
		"WNDSR": geo.Latlong{38.681808, -122.478747},
		"WEBRR": geo.Latlong{38.243881, -122.412142},
		"SPAMY": geo.Latlong{39.200661, -122.591042},
		"HUBRT": geo.Latlong{39.040228, -122.568314},
		"DRAXE": geo.Latlong{38.759, -122.389047},
		"BMBOO": geo.Latlong{38.892972, -122.233019},
		"RBUCL": geo.Latlong{39.070053, -122.02615},
		"GRTFL": geo.Latlong{38.35216944, -122.2314694},
		"TRUKN": geo.Latlong{37.71755833, -122.2145889},
		"DEDHD": geo.Latlong{38.33551666, -122.1128083},
		"HYPEE": geo.Latlong{37.88024444, -122.0674833},
		"COSMC": geo.Latlong{37.82606111, -122.0049},
		"TYDYE": geo.Latlong{37.689319, -122.268944},
		"ORRCA": geo.Latlong{38.610325, -121.551622},
		"MOGEE": geo.Latlong{38.336111, -121.389722},
		"TIPRE": geo.Latlong{38.205833, -121.035833},
		"SYRAH": geo.Latlong{37.99105, -121.103089},
		"RAIDR": geo.Latlong{38.0325, -122.5575},
		"CRESN": geo.Latlong{37.697475, -122.012019},
		"AAAME": geo.Latlong{37.770908, -122.082811},
		"ALLXX": geo.Latlong{37.729606, -122.064283},
		"HIRMO": geo.Latlong{37.92765, -122.14835},
		"CEXUR": geo.Latlong{37.934161, -122.252928},
		"WOULD": geo.Latlong{37.774508, -122.058064},
		"FINSH": geo.Latlong{37.651203, -122.257161},
		"HUSHH": geo.Latlong{37.7495, -122.338592},
		"AANET": geo.Latlong{38.530769, -122.497194},

		// Foster city
		"ROKME": geo.Latlong{37.5177778, -122.1191667},
		"DONGG": geo.Latlong{37.5891667, -122.2525000},
		"GUTTS": geo.Latlong{37.5552778, -122.1597222},
		// This GOBEC is wrong ... see the one below.
		// "GOBEC": geo.Latlong{37.5869444, -122.2547222},
		"WASOP": geo.Latlong{37.5391667, -122.1247222},
		"DUYET": geo.Latlong{37.5680556, -122.2547222},

		// DYAMD and YOSEM
		"ARCHI": geo.Latlong{37.490833,    -121.875500},
		"FRELY": geo.Latlong{37.510667,    -121.793167},
		"CEDES": geo.Latlong{37.550822,    -121.624586},
		"FLOWZ": geo.Latlong{37.592500,    -121.264833},
		"ALWYS": geo.Latlong{37.633500,    -120.959333},
		"LAANE": geo.Latlong{37.659000,    -120.747333},
		"DYAMD": geo.Latlong{37.699167,    -120.404500},

		"FAITH": geo.Latlong{37.401217,    -121.861900},
		"SOOIE": geo.Latlong{37.428500,    -121.607667},
		"FRIGG": geo.Latlong{37.465500,    -121.257333},
		"ZOMER": geo.Latlong{37.545333,    -120.631500},
		"SNORA": geo.Latlong{37.645500,    -119.806333},
		"YOSEM": geo.Latlong{37.762667,    -118.766667},

		// Final approaches into SFO (28L, 28R)
		"GOBEC": geo.Latlong{37.578833,    -122.252833},
		"JOSUF": geo.Latlong{37.592167,    -122.285500},
		"DARNE": geo.Latlong{37.593333,    -122.292333},
		"FABLA": geo.Latlong{37.597500,    -122.318833},
		"AXMUL": geo.Latlong{37.571500,    -122.257167},
		"WIBNI": geo.Latlong{37.516667,    -122.031333},
		"ANETE": geo.Latlong{37.463667,    -121.942667},
		"FATUS": geo.Latlong{37.486000,    -122.002333},
		"HEGOT": geo.Latlong{37.508000,    -122.061833},
		"MIUKE": geo.Latlong{37.552333,    -122.181167},
		"DIVEC": geo.Latlong{37.432833,    -121.935000},
		"CEPIN": geo.Latlong{37.536000,    -122.172833},
		"DUMBA": geo.Latlong{37.503500,    -122.096167},
		"GIRRR": geo.Latlong{37.495852,    -122.027167},
		"ZILED": geo.Latlong{37.495667,    -121.958167},

		// SJC arrivals (incl. reverse flow)
		"HITIR": geo.Latlong{37.323567, -122.007392},
		"JESEN": geo.Latlong{37.294831, -121.975569},
		"PUCKK": geo.Latlong{37.363500, -122.009667},
		
		// Personal entries
		"X_RSH": geo.Latlong{36.868582,  -121.691934},
		"X_BLH": geo.Latlong{37.2199471, -122.0425108},
		"X_HBR": geo.Latlong{37.309564,  -122.112378},
		"X_WSD": geo.Latlong{37.420995,  -122.268237}, // Woodside
		"X_PVY": geo.Latlong{37.38087,   -122.23319}, // Portola Valley
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

	// http://flightaware.com/resources/airport/SFO/STAR/SERFR+TWO+(RNAV)/pdf
	Serfr1 = geo.Procedure{
		Name: "SERFR2",
		Departure: false,
		Airport: "SFO",
		Waypoints: []geo.Waypoint{
			{"SERFR", geo.Latlong{},    0,     0,   0, false}, // Many aircraft skip SERFR
			{"NNRLI", geo.Latlong{}, 20000, 20000, 280, true},
			{"WWAVS", geo.Latlong{}, 15000, 19000, 280, true},
			{"EPICK", geo.Latlong{}, 10000, 15000, 280, true},
			{"EDDYY", geo.Latlong{},  6000,  6000, 240, true},  // Delay vectoring inside EPICK-EDDYY
			{"SWELS", geo.Latlong{},  4700,  4700, 240, false},
			{"MENLO", geo.Latlong{},  4000,  4000, 230, false},
		},
	}
)

func ListWaypoints() []string {
	ret := []string{}
	for k,_ := range KFixes { ret = append(ret,k) }
	sort.Strings(ret)
	return ret
}

// A version of this which looks up against SFO names
func FormValueNamedLatlong(r *http.Request, stem string) geo.NamedLatlong {
	vals := KFixes
	for k,v := range KAirports { vals[k] = v }
	return geo.FormValueNamedLatlong(r, vals, stem)
}


// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
