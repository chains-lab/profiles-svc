package handlers

import (
	"net/http"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/config"
)

type Service struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewService(cfg config.Config, log logium.Logger, a *app.App) Service {
	return Service{
		app: a,
		log: log,
		cfg: cfg,
	}
}

func (s Service) Log(r *http.Request) logium.Logger {
	return s.log
}
