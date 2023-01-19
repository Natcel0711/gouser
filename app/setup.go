package app

import (
	"errors"

	"github.com/Natcel0711/gouser/config"
	"github.com/Natcel0711/gouser/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupApp() error {
	err := config.LoadENV()
	if err != nil {
		return err
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	err = router.SetupRouters(r)
	if err != nil {
		return errors.New("error setting up router")
	}
	return nil
}
