package api

import (
	"SongLibraryApi/models"
	"SongLibraryApi/repositories"
	"SongLibraryApi/services"
	"SongLibraryApi/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// @Summary Get list of songs with filtering and pagination
// @Description Get a list of songs with optional filters and pagination
// @Tags songs
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song"
// @Param limit query int false "Limit results (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} models.Song
// @Failure 400 {object} map[string]string
// @Router /api/v1/songs [get]
func GetSongsHandler(c *gin.Context) {
	filter := make(map[string]string)
	if group := c.Query("group"); group != "" {
		filter["group_name"] = group
	}
	if song := c.Query("song"); song != "" {
		filter["song_title"] = song
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	repo := repositories.SongRepo{}
	songs, err := repo.GetSongs(filter, limit, offset)
	if err != nil {
		utils.Logger.Error("Ошибка при получении песен", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song verses with pagination
// @Description Get the text of a song split into verses with pagination
// @Tags songs
// @Param id path string true "Song ID"
// @Param limit query int false "Number of verses per page (default: 1)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} map[string][]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/songs/{id}/verses [get]
func GetSongVersesHandler(c *gin.Context) {
	id := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	repo := repositories.SongRepo{}
	text, err := repo.GetSongTextWithPagination(id, limit, offset)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warn("Песня не найдена", zap.String("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}

		utils.Logger.Error("Ошибка при получении текста песни", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verses": text})
}

// @Summary Add a new song to the library
// @Description Add a new song by making a request to the external API
// @Tags songs
// @Accept json
// @Produce json
// @Param input body models.Song true "Song details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/songs [post]
func AddSongHandler(c *gin.Context) {
	var input struct {
		Group string `json:"group" binding:"required"`
		Song  string `json:"song" binding:"required"`
	}

	// Проверяем корректность входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Logger.Warn("Неверный формат данных", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	repo := repositories.SongRepo{}
	exists, err := repo.CheckSongExists(input.Group, input.Song)
	if err != nil {
		utils.Logger.Error("Ошибка при проверке существования песни", zap.String("group", input.Group), zap.String("song", input.Song), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if exists {
		utils.Logger.Warn("Попытка создать существующую песню", zap.String("group", input.Group), zap.String("song", input.Song))
		c.JSON(http.StatusConflict, gin.H{"error": "Song already exists"})
		return
	}

	details, err := services.FetchSongDetails(input.Group, input.Song)
	if err != nil {
		utils.Logger.Error("Ошибка при запросе к внешнему API", zap.String("group", input.Group), zap.String("song", input.Song), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch song details"})
		return
	}

	releaseDate, err := time.Parse("02.01.2006", details.ReleaseDate)
	if err != nil {
		utils.Logger.Error("Ошибка при парсинге даты", zap.String("date", details.ReleaseDate), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse release date"})
		return
	}

	newSong := models.Song{
		GroupName:   input.Group,
		SongTitle:   input.Song,
		ReleaseDate: releaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	err = repo.CreateSong(newSong)
	if err != nil {
		utils.Logger.Error("Ошибка при создании песни", zap.String("group", input.Group), zap.String("song", input.Song), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	utils.Logger.Info("Песня успешно создана", zap.String("group", input.Group), zap.String("song", input.Song))
	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully"})
}

// @Summary Update an existing song
// @Description Update song details by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param input body models.Song true "Updated song details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/songs/{id} [put]
func UpdateSongHandler(c *gin.Context) {
	id := c.Param("id")

	repo := repositories.SongRepo{}
	var song models.Song
	if err := repo.FindSongByID(id, &song); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warn("Песня не найдена", zap.String("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}

		utils.Logger.Error("Ошибка при проверке существования песни", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Logger.Warn("Неверный формат данных", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := repo.UpdateSong(id, input); err != nil {
		utils.Logger.Error("Ошибка при обновлении песни", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	utils.Logger.Info("Песня успешно обновлена", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// @Summary Delete a song by ID
// @Description Delete an existing song by its ID
// @Tags songs
// @Param id path string true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/songs/{id} [delete]
func DeleteSongHandler(c *gin.Context) {
	id := c.Param("id")

	repo := repositories.SongRepo{}
	var song models.Song
	if err := repo.FindSongByID(id, &song); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warn("Песня не найдена", zap.String("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}

		utils.Logger.Error("Ошибка при проверке существования песни", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := repo.DeleteSong(id); err != nil {
		utils.Logger.Error("Ошибка при удалении песни", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	utils.Logger.Info("Песня успешно удалена", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
