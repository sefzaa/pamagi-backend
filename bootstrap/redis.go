package bootstrap

import (
	"context"
	"log"
	"fmt" // Tambahkan ini

	"github.com/redis/go-redis/v9"
)

func NewRedis(env *Env) *redis.Client {
	// Dapatkan alamat dari env.REDIS_HOST (pastikan struct Env kamu punya field ini)
	// Jika belum, gunakan fmt.Sprintf untuk menggabungkan host dan port
	addr := fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr, 
		Password: env.RedisPassword, 
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Gagal koneksi ke Redis: ", err)
	}

	log.Println("Berhasil terhubung ke Redis!")
	return client
}