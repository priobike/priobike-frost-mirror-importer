package things

// A location model from the SensorThings API.
type Location struct {
	Description  string   `json:"description"`
	EncodingType string   `json:"encodingType"`
	IotId        int      `json:"@iot.id"`
	Name     string `json:"name"`
	Location     struct { // GeoJSON
		Type     string `json:"type"`
		Geometry struct {
			Type        string        `json:"type"` // MultiLineString
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"location"`
}
