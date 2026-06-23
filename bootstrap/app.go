package bootstrap

import (
	
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Env *Env
	DB  *gorm.DB
	Redis *redis.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.DB = NewMySQLDatabase(app.Env)
	app.Redis = NewRedis(app.Env)
	return *app
}