package services

import (
	"SongLibraryApi/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func FetchSongDetails(group, song string) (*SongDetail, error) {
	config := utils.LoadConfig()
	url := fmt.Sprintf("%s?group=%s&song=%s", config.ExternalAPI, group, song)

	resp, err := http.Get(url)
	if err != nil {
		utils.Logger.Error("Не удалось выполнить запрос к внешнему API", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.Logger.Warn("Внешний API вернул ошибку", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Error("Не удалось прочитать ответ от внешнего API", zap.Error(err))
		return nil, err
	}

	var details SongDetail
	if err := json.Unmarshal(body, &details); err != nil {
		utils.Logger.Error("Не удалось распарсить ответ от внешнего API", zap.Error(err))
		return nil, err
	}

	return &details, nil
}
