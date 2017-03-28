package geo

import(
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"github.com/skypies/util/widget"
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

func FormValueNamedLatlong(r *http.Request, names map[string]Latlong, stem string) NamedLatlong {
	if name := strings.ToUpper(r.FormValue(stem+"_name")); name != "" {
		if _,exists := names[name]; !exists {
			return NamedLatlong{Name:"[UNKNOWN]"} //, fmt.Errorf("Waypoint '%s' not known", wp)
		}
		return NamedLatlong{name, names[name]}
	}

	return NamedLatlong{"", FormValueLatlong(r,stem)}
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


func (nl NamedLatlong)ToCGIArgs(stem string) string {
	v := url.Values{}
	widget.AddPrefixedValues(v, nl.Values(), stem)
	return v.Encode()
}

func (pos Latlong)Values() url.Values {
	v := url.Values{}
	v.Set("lat", fmt.Sprintf("%.5f", pos.Lat)) 
	v.Set("long", fmt.Sprintf("%.5f", pos.Long)) 
	return v
}

// you'll want widget.AddPrefixedValues(v, nl.Values(), "mystem")
func (nl NamedLatlong)Values() url.Values {
	v := url.Values{}
	v.Set("name", nl.Name)
	widget.AddValues(v, nl.Latlong.Values())
	return v
}

func (ln LatlongLine)Values() url.Values {
	v := url.Values{}
	widget.AddPrefixedValues(v, ln.From.Values(), "startpos")
	widget.AddPrefixedValues(v, ln.To.Values(), "endpos")
	return v
}
func (ln LatlongLine)ToCGIArgs(stem string) string {
	v := url.Values{}
	widget.AddPrefixedValues(v, ln.Values(), stem)
	return v.Encode()
}

// If fields absent or blank, returns {0.0, 0.0} -> {0.0, 0.0}
func FormValueLatlongLine(r *http.Request, stem string) LatlongLine {
	s := FormValueLatlong(r, stem+"_startpos")
	e := FormValueLatlong(r, stem+"_endpos")
	return s.LineTo(e)
}
