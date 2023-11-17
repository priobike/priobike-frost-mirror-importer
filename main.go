package main

import (
	"importer/filter"
	"importer/frost"
	"importer/things"
)

func main() {
	// Sync the things
	things.SyncThings()

	// Filter things
	filter.FilterThings()

	// Sync the frost data
	frost.Sync()
}
