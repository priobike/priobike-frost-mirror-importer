package frost

import (
	"bytes"
	"encoding/json"
	"importer/env"
	"importer/log"
	"importer/things"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func pruneOldData(client *http.Client) {
	log.Info.Println("Pruning things...")

	// Fetches pages and deletes things until no things are left
	for {
		pageUrl := env.SensorThingsProxyBaseUrl + "Things?"

		resp, err := client.Get(pageUrl)
		if err != nil {
			log.Warning.Println("Could not prune things:", err)
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Warning.Println("Could not prune things:", err)
			panic(err)
		}

		var thingsResponse struct {
			Value   []things.Thing `json:"value"`
			NextUri *string        `json:"@iot.nextLink"`
		}

		if err := json.Unmarshal(body, &thingsResponse); err != nil {
			log.Warning.Println("Could not prune things:", err)
			panic(err)
		}

		log.Info.Printf("Pruning %d things...", len(thingsResponse.Value))

		// Break if all things are deleted
		if len(thingsResponse.Value) == 0 {
			break
		}

		// Delete things
		for _, t := range thingsResponse.Value {
			req, err := http.NewRequest(http.MethodDelete, env.SensorThingsProxyBaseUrl+"Things("+strconv.Itoa(t.IotId)+")", nil)
			if err != nil {
				panic(err)
			}

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			// Panic if not 200.
			if !strings.Contains(resp.Status, "200") {
				panic(resp.Status)
			}
		}

	}

	log.Info.Println("Pruned things.")
}

func importThing(client *http.Client, thing things.Thing) {

	body, err := json.Marshal(thing)
	if err != nil {
		panic(err)
	}

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, env.SensorThingsProxyBaseUrl+"Things", bodyReader)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Panic if not 201
	if !strings.Contains(resp.Status, "201") {
		panic(resp.Status)
	}
}

func importNewData(client *http.Client) {
	log.Info.Println("Import things...")

	thingsList := []things.Thing{}

	things.Things.Range(func(key, value interface{}) bool {
		thingsList = append(thingsList, value.(things.Thing))
		return true
	})

	log.Info.Printf("Things list len: %d", len(thingsList))

	// Import all Things in Frost Server.
	for _, thing := range thingsList {
		importThing(client, thing)
	}

	log.Info.Println("Imported things.")
}

// Executes the sync
func Sync() {
	client := &http.Client{}

	// Currently not used
	// pruneOldData(client)

	importNewData(client)
}
