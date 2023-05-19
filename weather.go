package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

const apiUrl = "https://api.openweathermap.org/data/2.5/weather"

type Conditions struct {
	Summary            string
	TemperatureCelsius float64
}

func (c *Conditions) String() string {
	return fmt.Sprintf("%s %.1fºC", c.Summary, c.TemperatureCelsius)
}

type jsonResponse struct {
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func Current(location, token string) (Conditions, error) {
	url := FormatURL(location, token)

	resp, err := http.Get(url)
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close()

	conditions, err := ParseJSON(resp.Body)
	if err != nil {
		return Conditions{}, err
	}

	return conditions, nil
}

func FormatURL(location, token string) string {
	return fmt.Sprintf("%s?q=%s&appid=%s", apiUrl, location, token)
}

func ParseJSON(r io.Reader) (Conditions, error) {
	decoder := json.NewDecoder(r)

	response := jsonResponse{}
	if err := decoder.Decode(&response); err != nil {
		return Conditions{}, err
	}

	summary := response.Weather[0].Main
	celsius := convertKelvinToCelsius(response.Main.Temp)

	return Conditions{summary, celsius}, nil
}

func LocationFromArgs(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("location not provided")
	}
	locationArgs := args[1:]

	cityParts := []string{}
	countryCodeParts := []string{}
	for i := 0; i < len(locationArgs); i++ {
		cityParts = append(cityParts, locationArgs[i])
		if strings.Contains(locationArgs[i], ",") {
			countryCodeParts = append(countryCodeParts, locationArgs[i+1:]...)
			break
		}
	}

	city := strings.Join(cityParts, " ")
	countryCode := strings.Join(countryCodeParts, "")
	location := city + countryCode

	return location, nil
}

func convertKelvinToCelsius(kelvin float64) (celsius float64) {
	celsius = kelvin - 273.15
	return math.Round(celsius*100) / 100
}
