package frost

import (
	"bytes"
	"encoding/json"
	"importer/env"
	"importer/things"
	"net/http"
)

func pruneOldData(client *http.Client) {
	req, err := http.NewRequest(http.MethodDelete, env.SensorThingsProxyBaseUrl+"Things", nil)
	if err != nil {
		panic(err)
	}

	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func importNewData(client *http.Client) {
	thingsList := []things.Thing{}

	things.Things.Range(func(key, value interface{}) bool {
		thingsList = append(thingsList, value.(things.Thing))
		return true
	})

	body, err := json.Marshal(thingsList)
	if err != nil {
		panic(err)
	}

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, env.SensorThingsProxyBaseUrl+"Things", bodyReader)
	if err != nil {
		panic(err)
	}

	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func Sync() {
	client := &http.Client{}

	pruneOldData(client)

	importNewData(client)
}
