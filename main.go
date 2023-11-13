package main

import (
	"importer/frost"
	"importer/things"
)

func main() {
	// Sync the things.
	things.SyncThings()

	// Sync the frost data.
	frost.Sync()
}
