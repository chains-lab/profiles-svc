package handlers

import (
	"github.com/chains-lab/logium"
	"github.com/chains-lab/profiles-svc/internal/app"
	"github.com/chains-lab/profiles-svc/internal/config"
)

type Handler struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewHandler(cfg config.Config, log logium.Logger, a *app.App) Handler {
	return Handler{
		app: a,
		log: log,
		cfg: cfg,
	}
}
