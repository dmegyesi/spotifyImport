package tracks

import (
	"github.com/zmb3/spotify"
	log "github.com/sirupsen/logrus"
	"strings"
)

func SearchSong(client *spotify.Client, t Track, playList *PlayList, missingTracks *Tracks) bool {
	var found bool = false

	var artist string = t.Artist
	var title string = t.Title
	var album string = t.Album

	log.WithFields(log.Fields{
		"source": "json",
	}).Debugf("%s - %s (%s)", artist, title, album)

	var searchString string

	// if title has ' character, it breaks Spotify search, so let's remove it
	t.Title = strings.Replace(t.Title, "'", "", -1)

	if artist != "" && title != "" {
		searchString = "track:" + title + " artist:" + artist
	} else if album != "" && title != "" {
		searchString = "track:" + title + " album:" + album
	} else {
		// failsafe
		return false
	}

	// We take only the first result anyway, so don't waste resources
	// (we trust Spotify that the first result is the best/closest result)
	var limit int = 1

	results, err := client.SearchOpt(searchString, spotify.SearchTypeTrack, &spotify.Options{Limit: &limit})
	if err != nil {
		log.Fatal(err)
	}

	if results.Tracks.Total > 0 {

		log.WithFields(log.Fields{
			"source": "spotify",
		}).Debugf("%s - %s", results.Tracks.Tracks[0].Artists[0].Name, results.Tracks.Tracks[0].Name)

		log.WithFields(log.Fields{
			"source": "spotify",
		}).Infof("%s", results.Tracks.Tracks[0].ID)

		t.ID = results.Tracks.Tracks[0].ID

		playList.AddTrack(results.Tracks.Tracks[0].ID)

		found = true

	} else {

		log.WithFields(log.Fields{
			"source": "spotify",
		}).Errorf("Not found: %s - %s (%s)", artist, title, album)

		missingTracks.AddTrack(t)
	}

	return found
}
