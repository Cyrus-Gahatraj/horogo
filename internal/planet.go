package internal

type Planet struct {
	Name string
	Id   string
}

type Placement struct {
    Sign   string  `json:"sign"`
    Degree float64 `json:"degree"`
}

type Chart struct {
    Name string 					`json:"name"`
	Place string					`json:"place"`
    Ascendant Placement 			`json:"ascendant"`
    Planets map[string]Placement 	`json:"planets"`
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

