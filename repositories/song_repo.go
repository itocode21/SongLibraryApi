package repositories

import (
	"SongLibraryApi/models"
	"strings"
)

type SongRepo struct{}

func (r *SongRepo) GetSongs(filter map[string]string, limit, offset int) ([]models.Song, error) {
	var songs []models.Song
	query := DB.Model(&models.Song{})

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	query.Limit(limit).Offset(offset).Find(&songs)
	return songs, nil
}

func (r *SongRepo) CreateSong(song models.Song) error {
	return DB.Create(&song).Error
}

func (r *SongRepo) CheckSongExists(groupName, songTitle string) (bool, error) {
	var count int64
	err := DB.Model(&models.Song{}).Where("group_name = ? AND song_title = ?", groupName, songTitle).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SongRepo) UpdateSong(id string, song models.Song) error {
	return DB.Model(&models.Song{}).Where("id = ?", id).Updates(song).Error
}

func (r *SongRepo) FindSongByID(songID string, song *models.Song) error {
	return DB.First(song, songID).Error
}

func (r *SongRepo) DeleteSong(id string) error {
	return DB.Delete(&models.Song{}, id).Error
}

func (r *SongRepo) GetSongTextWithPagination(songID string, limit, offset int) ([]string, error) {
	var song models.Song
	if err := DB.First(&song, songID).Error; err != nil {
		return nil, err
	}

	verses := strings.Split(song.Text, "\n\n")

	totalVerses := len(verses)
	if offset >= totalVerses {
		return []string{}, nil
	}
	end := offset + limit
	if end > totalVerses {
		end = totalVerses
	}

	return verses[offset:end], nil
}
