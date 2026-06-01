package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const GEO_URL string = "https://photon.komoot.io/api"
const TIMEZONE_URL string = "https://timeapi.io/api/timezone/coordinate"


type geoResponse struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Name    string `json:"name"`
			Country string `json:"country"`
		} `json:"properties"`
	} `json:"features"`
}

type tzResponse struct {
	TimeZone         string `json:"timeZone"`
	CurrentUtcOffset struct {
		Seconds int `json:"seconds"`
	} `json:"currentUtcOffset"`
}

func GeocodePlace(place string) (lat, lon float64, err error) {
	params := url.Values{}
	params.Set("q", place)
	params.Set("limit", "1")
	reqURL := GEO_URL + "?" + params.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		return 0, 0, fmt.Errorf("geocode request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("reading response failed: %w", err)
	}

	var gr geoResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		return 0, 0, fmt.Errorf("parsing response failed: %w", err)
	}

	if len(gr.Features) == 0 {
		return 0, 0, fmt.Errorf("no results found for %q", place)
	}

	coords := gr.Features[0].Geometry.Coordinates
	if len(coords) < 2 {
		return 0, 0, fmt.Errorf("invalid coordinates for %q", place)
	}

	return coords[1], coords[0], nil
}

func GetTimezoneOffset(lat, lon float64) (int, error) {
	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	reqURL := TIMEZONE_URL + "?" + params.Encode()

	resp, err := http.Get(reqURL)
	if err != nil {
		return 0, fmt.Errorf("timezone request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response failed: %w", err)
	}

	var tr tzResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return 0, fmt.Errorf("parsing response failed: %w", err)
	}

	return tr.CurrentUtcOffset.Seconds, nil
}
