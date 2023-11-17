package things

import (
	"encoding/json"
	"fmt"
	"importer/env"
	"importer/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// A map that contains all things by their name.
var Things = &sync.Map{}

// Syncs things from specific pages
func syncThingsPage(page int) (more bool) {
	elementsPerPage := 100
	pageUrl := env.SensorThingsBaseUrl + "Things?" + url.QueryEscape(
		"$filter="+
			"Datastreams/properties/serviceName eq 'HH_STA_traffic_lights' "+
			"and (Datastreams/properties/layerName eq 'primary_signal' "+
			"  or Datastreams/properties/layerName eq 'signal_program' "+
			"  or Datastreams/properties/layerName eq 'cycle_second' "+
			"  or Datastreams/properties/layerName eq 'detector_car' "+
			"  or Datastreams/properties/layerName eq 'detector_bike') "+
			"&$expand=Datastreams,Locations"+
			"&$skip="+fmt.Sprintf("%d", page*elementsPerPage),
	)

	resp, err := http.Get(pageUrl)
	if err != nil {
		log.Warning.Println("Could not sync things:", err)
		panic(err)
	}
	defer resp.Body.Close()

	// Panic if not 200.
	if !strings.Contains(resp.Status, "200") {
		panic(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warning.Println("Could not sync things:", err)
		panic(err)
	}

	var thingsResponse struct {
		Value   []Thing `json:"value"`
		NextUri *string `json:"@iot.nextLink"`
	}
	if err := json.Unmarshal(body, &thingsResponse); err != nil {
		log.Warning.Println("Could not sync things:", err)
		panic(err)
	}

	for _, t := range thingsResponse.Value {
		// Add the thing to the things map.
		Things.Store(t.Name, t)
	}

	return thingsResponse.NextUri != nil
}

// Periodically sync the things from the SensorThings API.
func SyncThings() {
	log.Info.Println("Syncing things...")

	Things = &sync.Map{}

	// Fetch all pages of the SensorThings query.
	var page = 0
	for {
		// Make some parallel requests to speed things up.
		var wg sync.WaitGroup
		var foundMore = false
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(page int) {
				defer wg.Done()
				more := syncThingsPage(page)
				if more {
					foundMore = true
				}
			}(page)
			page++
		}
		log.Info.Printf("Bulk syncing things from pages %d-%d...", page-10, page-1)
		wg.Wait()
		if !foundMore {
			break
		}
	}

	log.Info.Println("Synced things.")
}
