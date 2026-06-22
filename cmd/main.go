package main

import (
	"pamagi/api/route" 
	"pamagi/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi Environment dan Database MySQL
	app := bootstrap.App()

	env := app.Env
	db := app.DB

	// Setup Gin Engine
	ginEngine := gin.Default()

	// Setup Route (Pastikan route.Setup diupdate untuk menerima *gorm.DB, bukan Mongo)
	route.Setup(env, db, ginEngine)

	// Jalankan server
	ginEngine.Run(":" + env.AppPort)
}