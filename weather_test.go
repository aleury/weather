package weather_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"weather"

	"github.com/google/go-cmp/cmp"
)

func TestFormatURL(t *testing.T) {
	t.Parallel()
	location := "London"
	token := "dummy_token"
	client := weather.NewClient(token)
	want := "https://api.openweathermap.org/data/2.5/weather?q=London&appid=dummy_token"
	got := client.FormatURL(location)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestParseJSON(t *testing.T) {
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

func TestLocationFromArgs(t *testing.T) {
	t.Parallel()
	type testCase struct {
		args []string
		want string
	}
	tests := []testCase{
		{args: []string{"/usr/bin/weather", "London"}, want: "London"},
		{args: []string{"/usr/bin/weather", "London,", "UK"}, want: "London,UK"},
		{args: []string{"/usr/bin/weather", "Berlin,", "DE"}, want: "Berlin,DE"},
		{args: []string{"/usr/bin/weather", "New", "York"}, want: "New York"},
		{args: []string{"/usr/bin/weather", "New", "York,", "US"}, want: "New York,US"},
		{args: []string{"/usr/bin/weather", "Kingston", "upon", "Hull,", "UK"}, want: "Kingston upon Hull,UK"},
	}
	for _, tc := range tests {
		got, err := weather.LocationFromArgs(tc.args)
		if err != nil {
			t.Fatalf("didn't expect an error: %s", err)
		}
		if !cmp.Equal(tc.want, got) {
			t.Errorf("want %q, got %q", tc.want, got)
		}
	}
}

func TestLocationFromArgsWithInvalidInput(t *testing.T) {
	t.Parallel()
	args := []string{"/usr/bin/weather"}
	_, err := weather.LocationFromArgs(args)
	if err == nil {
		t.Error("expected an error to be returned")
	}
}

func TestConditions(t *testing.T) {
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

func TestCurrent(t *testing.T) {
	t.Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("testdata/london.json")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	client := weather.NewClient("dummy_token")
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
