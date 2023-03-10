package router

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Natcel0711/gouser/database"
	"github.com/Natcel0711/gouser/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRouters(r *chi.Mux) error {
	db, err := database.StartDB()
	if err != nil {
		return errors.New("error connecting to DB: " + err.Error())
	}
	defer database.CloseDB()
	r.Get("/health", handlers.HealthCheck)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", handlers.GetAllUsers(db))
		r.Get("/{id}", handlers.GetUserByID(db))
		r.Get("/BySession/{sessionid}", handlers.GetUserBySession(db))
		r.Get("/ByEmail/{email}", handlers.GetUserByEmail(db))
		r.Post("/Session/", handlers.CreateSessionID(db))
		r.Post("/", handlers.CreateUser(db))
		r.Put("/", handlers.UpdateUser(db))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("HELLO, ENVs " + os.Getenv("TEST_ENV"))
	log.Println("listening on port: " + port)
	http.ListenAndServe("0.0.0.0:"+port, r)
	return nil
}
