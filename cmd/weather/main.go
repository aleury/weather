package main

import (
	"fmt"
	"os"

	"github.com/aleury/weather"
)

func main() {
	token := os.Getenv("OPENWEATHER_API_TOKEN")
	client, err := weather.NewOpenWeatherClient(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(weather.RunCLI(client))
}
