package models

type Config struct {
	Port       int    `default:"3333"`
	DBHost     string `default:"localhost"`
	DBPort     int    `default:"5432"`
	DBUser     string `default:"postgres"`
	DBPassword string `default:"example"`
	DBName     string `default:"postgres"`
}
