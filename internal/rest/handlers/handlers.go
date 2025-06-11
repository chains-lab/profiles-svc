package handlers

import (
	"github.com/chains-lab/profile-storage/internal/app"
	"github.com/chains-lab/profile-storage/internal/config"
	"github.com/chains-lab/profile-storage/internal/rest/presenters"
	"github.com/sirupsen/logrus"
)

type App interface {
}

type Presenter interface {
}

type Handlers struct {
	app       App
	presenter Presenter
	log       *logrus.Entry
	cfg       config.Config
}

func NewHandlers(cfg config.Config, log *logrus.Entry, app *app.App) Handlers {
	return Handlers{
		app:       app,
		cfg:       cfg,
		presenter: presenters.NewPresenter(log),
		log:       log,
	}
}
