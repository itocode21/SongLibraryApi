package models

import "time"

type Song struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	GroupName   string    `json:"group_name"`
	SongTitle   string    `json:"song_title"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}
