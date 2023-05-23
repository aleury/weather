[![Test](https://github.com/aleury/weather/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/aleury/weather/actions/workflows/test.yml)

`weather` is a command-line tool that makes it easy to fetch the current weather conditions from all around the world. 

Install it with:

```
go install github.com/aleury/weather/cmd/weather@latest
```

# What?

For example, run the following if you'd like to get the current weather conditions of London, UK:

```
$ weather London, UK
Cloudy 15.2ºC
```

## Getting started

This tool uses the [OpenWeather](https://openweathermap.org/api) API, so you'll need to sign in to the OpenWeather site to create an API token and add it your shell environment:

```
export OPENWEATHER_API_KEY=my_api_token
```

