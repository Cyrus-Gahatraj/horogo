package internal

type Planet struct {
	Name string
	Id   string
}

type Placement struct {
	Sign      string  `json:"sign"`
	Degree    float64 `json:"degree"`
	House     int     `json:"house"`
	Nakshatra string  `json:"nakshatra"`
	Pada      int     `json:"pada"`
}

type Chart struct {
	Name      string               `json:"name"`
	Gender    string               `json:"gender"`
	Place     string               `json:"place"`
	Ascendant Placement            `json:"ascendant"`
	Houses    [12]float64          `json:"houses"`
	Planets   map[string]Placement `json:"planets"`
}

var Planets = []Planet{
	{"Sun", "10"},
	{"Moon", "301"},
	{"Mercury", "199"},
	{"Venus", "299"},
	{"Mars", "499"},
	{"Jupiter", "599"},
	{"Saturn", "699"},
	{"Uranus", "799"},
	{"Neptune", "899"},
	{"Pluto", "999"},
}

var zodiacSigns = []string{
	"Aries", "Taurus", "Gemini", "Cancer",
	"Leo", "Virgo", "Libra", "Scorpio",
	"Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

var nakshatras = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira",
	"Ardra", "Punarvasu", "Pushya", "Ashlesha", "Magha",
	"Purva Phalguni", "Uttara Phalguni", "Hasta", "Chitra", "Swati",
	"Vishakha", "Anuradha", "Jyeshtha", "Mula", "Purva Ashadha",
	"Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
	"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
}

func calcNakshatra(lon float64) (string, int) {
	for lon < 0 {
		lon += 360
	}
	for lon >= 360 {
		lon -= 360
	}
	const span = 360.0 / 27.0
	const padaSpan = span / 4.0
	idx := int(lon / span)
	pada := int((lon-float64(idx)*span)/padaSpan) + 1
	return nakshatras[idx], pada
}

func calcHouse(lon, ascLon float64) int {
	diff := lon - ascLon
	if diff < 0 {
		diff += 360
	}
	return int(diff/30) + 1
}
