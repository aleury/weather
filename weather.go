// Package weather makes it easy to fetch the current weather conditions
// of cities from all around the world.
package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apiUrl = "https://api.openweathermap.org/data/2.5/weather"

type Client interface {
	Current(location string) (Conditions, error)
}

// OpenWeatherClient is an API client for fetching the current weather conditions
// from the Open Weather service.
type OpenWeatherClient struct {
	token      string
	BaseURL    string
	HttpClient *http.Client
}

// NewOpenWeatherClient returns a new instance of the Client configured with the given auth token.
func NewOpenWeatherClient(token string) (*OpenWeatherClient, error) {
	if token == "" {
		return nil, errors.New("missing api token")
	}
	return &OpenWeatherClient{
		token:      token,
		BaseURL:    apiUrl,
		HttpClient: &http.Client{},
	}, nil
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
	celsius = kelvin - 273.15
	return math.Round(celsius*100) / 100
}

// RunCLI runs the command-line interface using the given
// weather api client.
func RunCLI(client Client) int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "location not provided")
		return 1
	}
	location := strings.Join(os.Args[1:], " ")

	conditions, err := client.Current(location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't fetch weather conditions for location %q: %s\n", location, err)
		return 1
	}

	fmt.Println(conditions.String())
	return 0
}
