package geo

import(
	"fmt"
	"net/http"
	"strconv"
)

// cloned from util/widget
func formValueFloat64EatErrs(r *http.Request, name string) float64 {	
	if val,err := strconv.ParseFloat(r.FormValue(name), 64); err != nil {
		return 0.0
	} else {
		return val
	}
}
func formValueInt64EatErrs(r *http.Request, name string) int64 {
	val,_ := strconv.ParseInt(r.FormValue(name), 10, 64)
	return val
}

// Routines to read/write objects from CGI params

// If fields absent or blank, returns {0.0, 0.0}
func FormValueLatlong(r *http.Request, stem string) Latlong {
	lat  := formValueFloat64EatErrs(r, stem+"_lat")
	long := formValueFloat64EatErrs(r, stem+"_long")
	return Latlong{lat,long}
}

func (pos Latlong)ToCGIArgs(stem string) string {
	return fmt.Sprintf("%s_lat=%.5f&%s_long=%.5f", stem, pos.Lat, stem, pos.Long)
}

// If fields absent or blank, returns {0.0, 0.0} -> {0.0, 0.0}
func FormValueLatlongBox(r *http.Request, stem string) LatlongBox {
	return LatlongBox{
		SW: FormValueLatlong(r, stem+"_sw"),
		NE: FormValueLatlong(r, stem+"_ne"),
		Floor: formValueInt64EatErrs(r, stem+"_floor"),
		Ceil: formValueInt64EatErrs(r, stem+"ceil"),
	}
}

func (box LatlongBox)ToCGIArgs(stem string) string {
	str := fmt.Sprintf("%s&%s", box.SW.ToCGIArgs(stem+"_sw"), box.NE.ToCGIArgs(stem+"_ne"))
	if box.Floor > 0 { str += fmt.Sprintf("%s_floor=%d", stem, box.Floor) }
	if box.Ceil > 0 { str += fmt.Sprintf("%s_ceil=%d", stem, box.Ceil) }
	return str
}
