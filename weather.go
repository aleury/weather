package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

const apiUrl = "https://api.openweathermap.org/data/2.5/weather"

type Conditions struct {
	Summary            string
	TemperatureCelsius float64
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
	celcius := convertKevinToCelcius(response.Main.Temp)

	return Conditions{summary, celcius}, nil
}

func convertKevinToCelcius(kelvin float64) (celcius float64) {
	celcius = kelvin - 273.15
	return math.Round(celcius*100) / 100
}
