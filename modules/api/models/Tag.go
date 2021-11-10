package models

import "github.com/google/uuid"

type videoTag struct {
	Id    uuid.UUID `json:"id"`
	Video Video     `json:"video_id"`
	Name  *string   `json:"name"`
}

type playlistTag struct {
	Id       uuid.UUID    `json:"id"`
	Playlist UserPlaylist `json:"playlist_id"`
	Name     *string      `json:"name"`
}
