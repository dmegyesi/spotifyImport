package tracks

import (
	"github.com/zmb3/spotify"
)

type Tracks struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Artist string     `json:"artist"`
	Title  string     `json:"title"`
	Album  string     `json:"album"`
	ID     spotify.ID `json:"id"`
}

func (tracks *Tracks) AddTrack(t Track) []Track {
	tracks.Tracks = append(tracks.Tracks, t)
	return tracks.Tracks
}
