// Package weather makes it easy to fetch the current weather conditions
// of cities from all around the world.
package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apiUrl = "https://api.openweathermap.org/data/2.5/weather"

// OpenWeatherClient is an API client for fetching the current weather conditions
// from the Open Weather service.
type OpenWeatherClient struct {
	token      string
	BaseURL    string
	HttpClient *http.Client
}

// NewOpenWeatherClient returns a new instance of the Client configured with the given auth token.
func NewOpenWeatherClient(token string) *OpenWeatherClient {
	client := &OpenWeatherClient{
		token:      token,
		BaseURL:    apiUrl,
		HttpClient: &http.Client{},
	}
	return client
}

// FormatURL returns a URL for fetching the current weather of the given location.
func (c *OpenWeatherClient) FormatURL(location string) string {
	return fmt.Sprintf("%s?q=%s&appid=%s", c.BaseURL, url.QueryEscape(location), c.token)
}

// Current fetches the present weather conditions of the given location.
func (c *OpenWeatherClient) Current(location string) (Conditions, error) {
	url := c.FormatURL(location)

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Conditions{}, errors.New(resp.Status)
	}

	conditions, err := ParseJSON(resp.Body)
	if err != nil {
		return Conditions{}, err
	}

	return conditions, nil
}

// Conditions holds the weather summary and temperature of a particular location.
type Conditions struct {
	Summary            string
	TemperatureCelsius float64
}

// String formats the weather conditions as a string.
func (c Conditions) String() string {
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

// ParseJSON decodes the JSON response provided by the given io.Reader into
// the struct for weather conditions.
func ParseJSON(r io.Reader) (Conditions, error) {
	decoder := json.NewDecoder(r)

	response := jsonResponse{}
	if err := decoder.Decode(&response); err != nil {
		return Conditions{}, err
	}
	if len(response.Weather) < 1 {
		return Conditions{}, errors.New("invalid response: missing weather data")
	}

	summary := response.Weather[0].Main
	celsius := convertKelvinToCelsius(response.Main.Temp)

	return Conditions{summary, celsius}, nil
}

func convertKelvinToCelsius(kelvin float64) (celsius float64) {
	return kelvin - 273.15
}

func CurrentWeather(token, location string) (Conditions, error) {
	client := NewOpenWeatherClient(token)
	return client.Current(location)
}

// RunCLI runs the command-line interface using the given
// weather api client.
func RunCLI() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "location not provided")
		return 1
	}
	location := strings.Join(os.Args[1:], " ")

	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "missing api token")
		return 1
	}

	conditions, err := CurrentWeather(token, location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't fetch weather conditions for location %q: %s\n", location, err)
		return 1
	}
	fmt.Println(conditions)
	return 0
}
