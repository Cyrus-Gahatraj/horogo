package internal

type Planet struct {
	Name string
	Id   string
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

