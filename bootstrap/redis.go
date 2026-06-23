package bootstrap

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedis(env *Env) *redis.Client {
	// Koneksi ke host "redis" dan port 6379 (Sesuai dengan nama service di docker-compose)
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", 
		Password: "", // Default Redis di Docker kita tidak pakai password
		DB:       0,  // Default DB
	})

	// Test koneksi
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Gagal koneksi ke Redis: ", err)
	}

	log.Println("Berhasil terhubung ke Redis!")
	return client
}