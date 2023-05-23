//go:build integration
// +build integration

package weather_test

import (
	"github.com/aleury/weather"
	"github.com/rogpeppe/go-internal/testscript"
	"os"
	"testing"
)

func TestRunCLIScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
		Setup: func(env *testscript.Env) error {
			env.Setenv("OPENWEATHER_API_TOKEN", os.Getenv("OPENWEATHER_API_TOKEN"))
			return nil
		},
	})
}

func TestOpenWeatherMapAPIReturnsCurrentWeatherConditions(t *testing.T) {
	t.Parallel()
	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		t.Skip("Please set a valid API key in the environment variable OPENWEATHER_API_TOKEN")
	}
	client := weather.NewClient(token)
	cond, err := client.Current("London")
	if err != nil {
		t.Fatal(err)
	}
	if cond.Summary == "" {
		t.Errorf("empty summary: %#v", cond)
	}
}
