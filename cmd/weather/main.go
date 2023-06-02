package main

import (
	"os"

	"github.com/aleury/weather"
)

func main() {
	os.Exit(weather.RunCLI())
}
