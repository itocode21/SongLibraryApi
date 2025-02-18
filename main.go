package main

import (
	"SongLibraryApi/api"
	"SongLibraryApi/repositories"
	"SongLibraryApi/utils"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "SongLibraryApi/docs" // Подключение документации Swagger

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config := utils.LoadConfig()
	utils.InitLogger(config.LogLevel)
	repositories.ConnectDB()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(utils.RequestLogger())
	api.SetupRoutes(r)

	port := config.Port
	utils.Logger.Info("Запуск сервера", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		utils.Logger.Fatal("Ошибка при запуске сервера", zap.Error(err))
		log.Fatalf("Failed to start server: %v", err)
	}
}
