package bootstrap

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLDatabase(env *Env) *gorm.DB {
	// Format DSN: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Menampilkan log query di terminal saat development
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Berhasil terkoneksi ke MySQL menggunakan GORM!")
	return db
}