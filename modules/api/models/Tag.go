package models

import "github.com/google/uuid"

type VideoTag struct {
	Id    uuid.UUID `json:"id"`
	Video Video     `json:"video_uuid"`
	Name  *string   `json:"name"`
}

type UserPlaylistTag struct {
	Id       uuid.UUID    `json:"id"`
	Playlist UserPlaylist `json:"playlist_uuid"`
	Name     *string      `json:"name"`
}

type UserTag struct {
	Id       uuid.UUID `json:"id"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Name     *string   `json:"name"`
}
