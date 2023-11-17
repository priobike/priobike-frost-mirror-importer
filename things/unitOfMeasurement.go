package things

// A traffic light datastream unit of measurement model from the SensorThings API.
type UnitOfMeasurement struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Definition string `json:"definition"`
}
