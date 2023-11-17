package things

// A traffic light datastream model from the SensorThings API.
type Datastream struct {
	IotId             int               `json:"@iot.id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	ObservationType   string            `json:"observationType"`
	UnitOfMeasurement UnitOfMeasurement `json:"unitOfMeasurement"`
	ObservedProperty  ObservedProperty  `json:"ObservedProperty"`
	Sensor            Sensor            `json:"Sensor"`
	Properties        struct {
		LayerName string `json:"layerName"`
	} `json:"properties"`
}
