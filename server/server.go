package server

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"transactions/crud"
	"transactions/models"
	"transactions/stores"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
)

type Server struct {
	Router *chi.Mux
	db     *sql.DB
	port   int
}

func NewServer() *Server {
	var c models.Config
	err := envconfig.Process("TRANSACTIONS", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err := getDB(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/account", crud.NewAccountCrud(*stores.NewAccountStore(db)).Routes())
	router.Mount("/transaction", crud.NewTransactionCrud(*stores.NewTransactionStore(db)).Routes())

	return &Server{
		Router: router,
		db:     db,
		port:   c.Port,
	}
}

func (s *Server) Start() {
	log.Printf("Server running on port %d", s.port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.port), s.Router))
}

func getDB(c models.Config) (*sql.DB, error) {
	db, err := stores.GetDB(
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
