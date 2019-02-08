package main

import (
	"encoding/json"
	"github.com/dmegyesi/spotifyImport/tracks"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type spotifyCredentials struct {
	SPOTIFY_ID     string `json:"SPOTIFY_ID"`
	SPOTIFY_SECRET string `json:"SPOTIFY_SECRET"`
}

func importTracksFromFile(destination *tracks.Tracks) {
	jsonFile, err := os.Open("tracks.json")

	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Successfully opened file")
	defer jsonFile.Close()

	// Import the entire array of tracks in one batch, only iterate later
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &destination)

	log.Infof("Total songs imported: %d songs", len(destination.Tracks))
}

func parseSpotifyCredentials(fileName string) spotifyCredentials {
	jsonFile, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Import the entire array of tracks in one batch, only iterate later
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var creds spotifyCredentials
	json.Unmarshal(byteValue, &creds)

	return creds
}
