package main

import (
	"encoding/json"
	"fmt"
	"github.com/dmegyesi/spotifyImport/tracks"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"net/http"
	"strings"
)

// redirectID is the OAuth redirect ID for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectID = "http://localhost:8080/callback"

const LogLevel = log.DebugLevel //log.InfoLevel

const credentialsFile = "account.json"

var (
	auth  = spotify.NewAuthenticator(redirectID, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPrivate)
	ch    = make(chan *spotify.Client)
	state = "ignore-audience-pointless-actress-spinach-manhunt"
)

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

func main() {

	log.SetLevel(LogLevel)

	credentials := parseSpotifyCredentials(credentialsFile)

	auth.SetAuthInfo(credentials.SPOTIFY_ID, credentials.SPOTIFY_SECRET)

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// if we receive a rate limit from Spotify, don't give up
	client.AutoRetry = true

	var parsedTracks tracks.Tracks
	var missingTracks tracks.Tracks
	var playList tracks.PlayList

	importTracksFromFile(&parsedTracks)

	// Search songs with the Spotify API
	for i, track := range parsedTracks.Tracks {
		tracks.SearchSong(client, track, &playList, &missingTracks)

		// If we're debugging, stop after 200 records
		if LogLevel == log.DebugLevel {
			if i == 200 {
				break
			}
		}
	}

	log.Infof("Parse errors: %d songs", len(missingTracks.Tracks))

	// For the missing tracks, try to do data cleaning and search again

	var notFoundTracks int = 0

	for _, track := range missingTracks.Tracks {
		var artist string = track.Artist
		var title string = track.Title
		var album string = track.Album

		// If the title or artist name contains a ( character, drop everything after it
		if strings.Contains(title, "(") {
			track.Title = strings.Split(title, "(")[0]
			if strings.Contains(artist, "(") {
				track.Artist = strings.Split(artist, "(")[0]
			}
		}

		found := tracks.SearchSong(client, track, &playList, &missingTracks)
		if !found {
			log.WithFields(log.Fields{
				"source": "spotify",
			}).Errorf("I gave up on: %s - %s (%s)", track.Artist, track.Title, album)

			notFoundTracks++
		}
	}

	log.Infof("Not found: %d songs", notFoundTracks)

	if LogLevel == log.DebugLevel {
		missingList, _ := json.Marshal(missingTracks)
		log.Debugln(string(missingList))
	}

	// Save songs to a new playlist
	tracks.UploadPlaylist(client, &playList)

}
