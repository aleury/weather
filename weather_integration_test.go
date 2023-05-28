//go:build integration
// +build integration

package weather_test

import (
	"github.com/aleury/weather"
	"os"
	"testing"
)

func TestOpenWeatherAPI_ReturnsCurrentWeatherConditions(t *testing.T) {
	t.Parallel()
	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		t.Skip("Please set a valid API key in the environment variable OPENWEATHER_API_TOKEN")
	}
	client, err := weather.NewOpenWeatherClient(token)
	if err != nil {
		t.Fatal(err)
	}
	cond, err := client.Current("London")
	if err != nil {
		t.Fatal(err)
	}
	if cond.Summary == "" {
		t.Errorf("empty summary: %#v", cond)
	}
}
