package models

type Config struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string `default:"postgres"`
	Password string `default:"example"`
	DBName   string `default:"postgres"`
}
