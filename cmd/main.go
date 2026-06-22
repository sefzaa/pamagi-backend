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

// Anotasi ini akan dibaca oleh Swagger
// @title Pamagi API
// @version 1.0
// @description Ini adalah dokumentasi API untuk backend Pamagi.
// @host localhost:8080
// @BasePath /
func main() {
	// Inisialisasi Environment dan Database MySQL
	app := bootstrap.App()

	env := app.Env
	db := app.DB

	// Setup Gin Engine
	ginEngine := gin.Default()

	// Setup Route Utama
	route.Setup(env, db, ginEngine)

	// Setup Route khusus untuk UI Swagger
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Jalankan server
	ginEngine.Run(":" + env.AppPort)
}