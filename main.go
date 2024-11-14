package main

import (
	"database/sql"
	"log"
	"transactions/models"
	"transactions/stores"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

func main() {
	var c models.Config
	err := envconfig.Process("TRANSACTIONS", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := stores.GetDB(
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	err = ping(db)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ping(db *sql.DB) error {
	return db.Ping()
}
