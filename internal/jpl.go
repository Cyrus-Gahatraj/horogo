package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const FORMAT string = "json"
const JPL_URL string = "https://ssd.jpl.nasa.gov/api/horizons.api"

func getFormattedTime(t time.Time) string {
	return t.Format("2006-01-02T15:04")
}

func eclipticToZodiac(lon float64) (string, float64) {
	for lon < 0 {
		lon += 360
	}
	for lon >= 360 {
		lon -= 360
	}
	idx := int(lon) / 30
	deg := lon - float64(idx*30)
	return zodiacSigns[idx], deg
}

const obliquity = 23.439281

func julianDay(year, month, day, hour, minute, second int) float64 {
	y, m := year, month
	if m <= 2 {
		y--
		m += 12
	}
	A := math.Floor(float64(y) / 100)
	B := 2 - A + math.Floor(A/4)
	jdn := math.Floor(365.25*float64(y+4716)) +
		math.Floor(30.6001*float64(m+1)) +
		float64(day) + B - 1524.5
	return jdn + float64(hour-12)/24 + float64(minute)/(24*60) + float64(second)/(24*3600)
}

func gmst(jd float64) float64 {
	return math.Mod(280.46061837+360.98564736629*(jd-2451545.0), 360)
}

func ascendant(ramcDeg, latDeg float64) float64 {
	ramcRad := ramcDeg * math.Pi / 180
	latRad := latDeg * math.Pi / 180
	epsRad := obliquity * math.Pi / 180

	ascRad := math.Atan2(-math.Cos(ramcRad),
		math.Sin(epsRad)*math.Tan(latRad)+
			math.Cos(epsRad)*math.Sin(ramcRad))
	ascDeg := ascRad * 180 / math.Pi
	return math.Mod(ascDeg+360, 360)
}

func (person Person) GetPlanetryPosition() Chart {
	tz := time.FixedZone("Local", person.TZOffset)
	localTime := time.Date(person.Year, time.Month(person.Month), person.Day,
		person.Hour, person.Minute, person.Second, 0, tz)
	utcTime := localTime.UTC()

	jd := julianDay(utcTime.Year(), int(utcTime.Month()), utcTime.Day(),
		utcTime.Hour(), utcTime.Minute(), utcTime.Second())
	gmstDeg := gmst(jd)
	ramcDeg := math.Mod(gmstDeg+person.Lon, 360)
	ascLon := ascendant(ramcDeg, person.Lat)
	ascSign, ascDeg := eclipticToZodiac(ascLon)
	fmt.Printf("Asc %s %.1f\n", ascSign, ascDeg)

	startTime := utcTime
	endTime := startTime.Add(time.Hour)

	params := url.Values{}
	params.Set("format", FORMAT)
	params.Set("START_TIME", getFormattedTime(startTime))
	params.Set("STOP_TIME", getFormattedTime(endTime))
	params.Set("STEP_SIZE", "1")
	params.Set("CSV_FORMAT", "YES")
	params.Set("QUANTITIES", "31")

	chart := Chart{
		Name: person.Name,
		Ascendant: Placement{
			Sign: ascSign,
			Degree: ascDeg,
		},
		Planets: map[string]Placement{},
	}


	for _, planet := range Planets {
		params.Set("COMMAND", planet.Id)
		wholeURL := JPL_URL + "?" + params.Encode()

		resp, err := http.Get(wholeURL)
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}

		var apiResp struct {
			Result string `json:"result"`
		}
		if err := json.Unmarshal(body, &apiResp); err != nil {
			panic(err)
		}

		lon := parseEclipticLongitude(apiResp.Result)
		sign, degree := eclipticToZodiac(lon)
		chart.Planets[planet.Name] = Placement{
			Sign: sign,
			Degree: degree,
		}
		fmt.Printf("%s %s %.1f\n", planet.Name, sign, degree)
	}

	return chart
}

func parseEclipticLongitude(result string) float64 {
	lines := strings.Split(result, "\n")
	inData := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "$$SOE" {
			inData = true
			continue
		}
		if trimmed == "$$EOE" {
			break
		}
		if !inData || trimmed == "" {
			continue
		}
		r := csv.NewReader(strings.NewReader(line))
		records, err := r.Read()
		if err != nil || len(records) < 4 {
			continue
		}
		lonStr := strings.TrimSpace(records[3])
		if lonStr == "" || lonStr == "n.a." {
			continue
		}
		var lon float64
		if _, err := fmt.Sscanf(lonStr, "%f", &lon); err != nil {
			continue
		}
		return lon
	}
	return 0
}

