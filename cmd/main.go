package main

import (
	"os"
	"weather"
)

func main() {
	os.Exit(weather.RunCLI())
}
