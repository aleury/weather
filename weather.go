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

type Client struct {
	token      string
	BaseURL    string
	HttpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		BaseURL:    apiUrl,
		HttpClient: &http.Client{},
	}
}

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

func (c *Client) FormatURL(location string) string {
	return fmt.Sprintf("%s?q=%s&appid=%s", c.BaseURL, location, c.token)
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
