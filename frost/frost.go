package frost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"importer/env"
	"importer/log"
	"importer/things"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func pruneThingsPage(client *http.Client, page int) (more bool) {
	elementsPerPage := 100
	pageUrl := env.SensorThingsProxyBaseUrl + "Things?" + url.QueryEscape(
		"$skip="+fmt.Sprintf("%d", page*elementsPerPage),
	)

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

	for _, t := range thingsResponse.Value {
		// Add the thing to the things map.
		req, err := http.NewRequest(http.MethodDelete, env.SensorThingsProxyBaseUrl+"Things("+strconv.Itoa(t.IotId)+")", nil)
		if err != nil {
			panic(err)
		}

		_, err = client.Do(req)
		if err != nil {
			panic(err)
		}
	}

	return thingsResponse.NextUri != nil
}

func pruneOldData() {
	log.Info.Println("Pruning things...")

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
				more := pruneThingsPage(&http.Client{Timeout: time.Duration(10) * time.Second}, page)
				if more {
					foundMore = true
				}
			}(page)
			page++
		}
		log.Info.Printf("Bulk pruning things from pages %d-%d...", page-10, page-1)
		wg.Wait()
		if !foundMore {
			break
		}
	}

	log.Info.Println("Pruned things.")
}

func importNewData(client *http.Client) {
	log.Info.Println("Import things...")

	thingsList := []things.Thing{}

	things.Things.Range(func(key, value interface{}) bool {
		thingsList = append(thingsList, value.(things.Thing))
		return true
	})

	// for _, thing := range thingsList {
	thing := thingsList[0]
	
	body, err := json.Marshal(thing)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(http.MethodPost, env.SensorThingsProxyBaseUrl+"Things", bodyReader)
	if err != nil {
		panic(err)
	}

	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	log.Info.Println("Thing " + strconv.Itoa(thing.IotId) + " imported.")
	// }

	log.Info.Println("Imported things.")
}

func Sync() {
	client := &http.Client{}

	pruneOldData()

	importNewData(client)
}
