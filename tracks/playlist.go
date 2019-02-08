package tracks

import (
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"time"
)

type PlayList struct {
	Ids []spotify.ID `json:"ids"`
}

func (playlist *PlayList) AddTrack(item spotify.ID) []spotify.ID {
	playlist.Ids = append(playlist.Ids, item)
	return playlist.Ids
}

func UploadPlaylist(client *spotify.Client, playList *PlayList) {

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	playlistName := "Google Music Import " + time.Now().Format(time.RFC850)
	spotifyPlaylist, err := client.CreatePlaylistForUser(user.ID, playlistName, "My imported playlist from "+
		"Google Music", false)

	if err != nil {
		log.Fatal(err)
	}

	// Spotify only allows up to 100 tracks in one request
	var i int = 0

	for i < len(playList.Ids) {
		var amount int

		if 100 < len(playList.Ids[i:]) {
			amount = 100
		} else {
			amount = len(playList.Ids[i:])
		}

		snapshot, err := client.AddTracksToPlaylist(spotifyPlaylist.ID, playList.Ids[i:i+amount]...)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Playlist snapshot ID: %s", snapshot)

		i += amount
	}

	log.Infof("New playlist created with %d songs.", i)
}
