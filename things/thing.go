package things

// A traffic light thing from the SensorThings API.
type Thing struct {
	IotId       int    `json:"@iot.id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Properties  struct {
		LaneType        string `json:"laneType"`
		TrafficLightsId string `json:"trafficLightsId"` // This is the crossing.
	} `json:"properties"`
	Datastreams []Datastream `json:"Datastreams"`
	Locations   []Location   `json:"Locations"`
}
