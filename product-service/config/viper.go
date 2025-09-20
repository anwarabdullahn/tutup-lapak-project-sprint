package config

import "github.com/spf13/viper"

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	_ = v.ReadInConfig() // ignore if missing
	v.AutomaticEnv()

	// defaults
	v.SetDefault("PORT", 3003)
	v.SetDefault("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	return v
}
