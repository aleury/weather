package weather_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aleury/weather"

	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

type testClient struct{}

func (c *testClient) Current(location string) (weather.Conditions, error) {
	if location != "London" {
		return weather.Conditions{}, errors.New("unknown location")
	}
	conditions := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	return conditions, nil
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"weather": func() int {
			client := testClient{}
			return weather.RunCLI(&client)
		},
	}))
}
func TestRunCLIScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}

func TestFormatURL_ReturnsURLWithLocationAndToken(t *testing.T) {
	t.Parallel()
	type testCase struct {
		location string
		token    string
		want     string
	}
	tests := []testCase{
		{location: "London", token: "dummy_token", want: "https://api.openweathermap.org/data/2.5/weather?q=London&appid=dummy_token"},
		{location: "London,UK", token: "dummy_token", want: "https://api.openweathermap.org/data/2.5/weather?q=London%2CUK&appid=dummy_token"},
		{location: "New York City", token: "dummy_token", want: "https://api.openweathermap.org/data/2.5/weather?q=New+York+City&appid=dummy_token"},
		{location: "New York City,US", token: "dummy_token", want: "https://api.openweathermap.org/data/2.5/weather?q=New+York+City%2CUS&appid=dummy_token"},
	}
	for _, tc := range tests {
		client, err := weather.NewOpenWeatherClient(tc.token)
		if err != nil {
			t.Fatalf("got an unexpected error %s", err)
		}
		got := client.FormatURL(tc.location)
		if tc.want != got {
			t.Errorf("want %q, got %q", tc.want, got)
		}
	}
}

func TestParseJSON_ReturnsWeatherConditonsForValidInput(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/london.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	want := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	got, err := weather.ParseJSON(f)
	if err != nil {
		t.Fatal("didn't expect an error parsing json")
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseJSON_ReturnsErrorForInvalidResponse(t *testing.T) {
	t.Parallel()
	response := bytes.NewBuffer([]byte(`{"error": "something went wrong"}`))
	_, err := weather.ParseJSON(response)
	if err == nil {
		t.Error("expected to get an error")
	}
}

func TestConditions_CanBeFormattedAsAString(t *testing.T) {
	conditions := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	want := "Drizzle 7.2ºC"
	got := conditions.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestNewOpenWeatherClient_ReturnsErrorForInvalidToken(t *testing.T) {
	want := "missing api token"
	_, err := weather.NewOpenWeatherClient("")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !cmp.Equal(err.Error(), want) {
		t.Error(cmp.Diff(err.Error(), want))
	}
}

func TestCurrent_ReturnsPresentWeatherConditonsOfAValidLocation(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("testdata/london.json")
		if err != nil {
			t.Fatalf("reading test data: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	client, err := weather.NewOpenWeatherClient("dummy_token")
	if err != nil {
		t.Fatalf("got an unexpected error %s", err)
	}
	client.BaseURL = server.URL
	client.HttpClient = server.Client()

	want := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	got, err := client.Current("London,UK")
	if err != nil {
		t.Fatalf("didn't expect an error")
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCurrent_ReturnsErrorWhenAPIRespondsWithUnknownLocation(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "unknown location"}`))
	}))
	defer server.Close()

	client, err := weather.NewOpenWeatherClient("dummy_token")
	if err != nil {
		t.Fatalf("got an unexpected error %s", err)
	}
	client.BaseURL = server.URL
	client.HttpClient = server.Client()

	_, err = client.Current("unknown")
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

func ExampleOpenWeatherClient_FormatURL() {
	client, err := weather.NewOpenWeatherClient("dummy_token")
	if err != nil {
		panic(err)
	}
	url := client.FormatURL("London,UK")
	fmt.Println(url)
	// Output:
	// https://api.openweathermap.org/data/2.5/weather?q=London%2CUK&appid=dummy_token
}

func ExampleParseJSON() {
	f, err := os.Open("testdata/london.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	conditions, err := weather.ParseJSON(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(conditions.String())
	// Output:
	// Drizzle 7.2ºC
}

func ExampleConditions_String() {
	conditions := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	fmt.Println(conditions.String())
	// Output:
	// Drizzle 7.2ºC
}
