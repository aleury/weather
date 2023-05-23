package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

const apiUrl = "https://api.openweathermap.org/data/2.5/weather"

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

// Client is an API client for fetching the current weather conditions
// from the Open Weather Map service.
type Client struct {
	token      string
	BaseURL    string
	HttpClient *http.Client
}

// NewClient returns a new instance of the Client configured with the given auth token.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		BaseURL:    apiUrl,
		HttpClient: &http.Client{},
	}
}

// RunCLI runs the command-line interface for weather.
func RunCLI() {
	location, err := LocationFromArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		log.Fatal("missing open weather api token")
	}
	client := NewClient(token)

	conditions, err := client.Current(location)
	if err != nil {
		log.Fatalf("error fetching current weather conditions: %s\n", err)
	}

	fmt.Println(conditions.String())
}

// Current fetches the present weather conditions of the given location.
func (c *Client) Current(location string) (Conditions, error) {
	url := c.FormatURL(location)

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf("request failed: %s", resp.Status)
	}

	conditions, err := ParseJSON(resp.Body)
	if err != nil {
		return Conditions{}, err
	}

	return conditions, nil
}

// FormatURL returns a URL for fetching the current weather of the given location.
func (c *Client) FormatURL(location string) string {
	return fmt.Sprintf("%s?q=%s&appid=%s", c.BaseURL, location, c.token)
}

// ParseJSON decodes the JSON response provided by the given io.Reader into
// the struct for weather conditions.
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

// LocationFromArgs parses the location from the arguments given to
// the command-line interface.
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
