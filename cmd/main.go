package main

import (
	"pamagi/api/route"
	"pamagi/bootstrap"

	"github.com/gin-gonic/gin"

	// Import library swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Import folder docs hasil generate (PENTING!)
	_ "pamagi/docs"
)

// @title Pamagi API
// @version 1.0
// @description Ini adalah dokumentasi API untuk backend Pamagi.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host pamagi.mydm.cloud
// @BasePath /
func main() {
	// Inisialisasi Environment dan Database MySQL
	app := bootstrap.App()

	env := app.Env
	db := app.DB
	redis := app.Redis
	// Setup Gin Engine
	ginEngine := gin.Default()

	// Setup Route Utama
	route.Setup(env, db, redis, ginEngine)

	// Setup Route khusus untuk UI Swagger
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Jalankan server
	ginEngine.Run(":" + env.AppPort)
}