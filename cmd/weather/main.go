package main

import (
	"github.com/aleury/weather"
	"os"
)

func main() {
	os.Exit(weather.RunCLI())
}
