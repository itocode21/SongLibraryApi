package repositories

import (
	"SongLibraryApi/models"
	"SongLibraryApi/utils"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*sql.DB, func()) {
	config := utils.LoadConfig()

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, "music_library", config.DBPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		_ = DB.Exec("DROP TABLE IF EXISTS songs CASCADE")
		_ = db.Close()
	}

	return db, cleanup
}

func TestCreateTable(t *testing.T) {

	db, cleanup := setupTestDB()
	defer cleanup()

	var tableName string
	err := db.QueryRow(`SELECT table_name FROM information_schema.tables WHERE table_name = 'songs'`).Scan(&tableName)
	if err == nil {
		t.Fatalf("Таблица 'songs' уже существует перед созданием")
	}

	err = DB.AutoMigrate(&models.Song{})
	if err != nil {
		t.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	err = db.QueryRow(`SELECT table_name FROM information_schema.tables WHERE table_name = 'songs'`).Scan(&tableName)
	if err != nil {
		t.Fatalf("Таблица 'songs' не была создана")
	}
	if tableName != "songs" {
		t.Fatalf("Создана неверная таблица. Ожидалась 'songs', получена '%s'", tableName)
	}
}

func TestDropTable(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	_, err := db.Exec(`
        CREATE TABLE songs (
            id SERIAL PRIMARY KEY,
            group_name VARCHAR(255) NOT NULL,
            song_title VARCHAR(255) NOT NULL,
            release_date DATE,
            text TEXT,
            link VARCHAR(255)
        );
    `)
	assert.NoError(t, err)

	_, err = db.Exec("DROP TABLE songs;")
	assert.NoError(t, err)

	var tableName string
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_name = 'songs'").Scan(&tableName)
	assert.Error(t, err)
}
func TestCRUDOperations(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()

	_, err := db.Exec(`
        CREATE TABLE songs (
            id SERIAL PRIMARY KEY,
            group_name VARCHAR(255) NOT NULL,
            song_title VARCHAR(255) NOT NULL,
            release_date DATE,
            text TEXT,
            link VARCHAR(255)
        );
    `)
	assert.NoError(t, err)

	_, err = db.Exec(`
        INSERT INTO songs (group_name, song_title, release_date, text, link)
        VALUES ('Muse', 'Supermassive Black Hole', '2006-07-16', 'Ooh baby...', 'https://youtube.com/...');
    `)
	assert.NoError(t, err)

	var group, title, text, link string
	var releaseDate time.Time
	err = db.QueryRow(`SELECT group_name, song_title, release_date, text, link FROM songs WHERE id = 1`).Scan(&group, &title, &releaseDate, &text, &link)
	assert.NoError(t, err)

	formattedReleaseDate := releaseDate.Format("2006-01-02")

	assert.Equal(t, "Muse", group)
	assert.Equal(t, "Supermassive Black Hole", title)
	assert.Equal(t, "2006-07-16", formattedReleaseDate)
	assert.Contains(t, text, "Ooh baby...")
	assert.Equal(t, "https://youtube.com/...", link)
}
