package weather

// Temperature is a type that stores its value with the unit Kelvin.
type Temperature float64

// Celsius returns the value of the temperature in Celsius.
func (t Temperature) Celsius() float64 {
	return float64(t) - 273.15
}
