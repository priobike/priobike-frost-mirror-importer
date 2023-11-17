package things

type Sensor struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	EncodingType string `json:"encodingType"`
	Metadata     string `json:"metadata"`
}
