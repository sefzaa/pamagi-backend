package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	AppPort                string `mapstructure:"APP_PORT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASSWORD"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenSecret      string `mapstructure:"JWT_SECRET"`
	AccessTokenExpiryMin   int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	RefreshTokenExpiryDays int    `mapstructure:"REFRESH_TOKEN_EXPIRE_DAYS"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	GmailSenderEmail string `mapstructure:"GMAIL_SENDER_EMAIL"`
	GmailAppPassword string `mapstructure:"GMAIL_APP_PASSWORD"`
	// BrevoAPIKey            string `mapstructure:"BREVO_API_KEY"`
	// BrevoSenderEmail       string `mapstructure:"BREVO_SENDER_EMAIL"`
	// ResendAPIKey           string `mapstructure:"RESEND_API_KEY"`
	// ResendSenderEmail      string `mapstructure:"RESEND_SENDER_EMAIL"`

}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}