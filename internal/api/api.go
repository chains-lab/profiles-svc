package api

import (
	"context"
	"net/http"

	"github.com/chains-lab/profile-storage/internal/app"
	"github.com/chains-lab/profile-storage/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Api struct {
	server *http.Server
	//router   *chi.Mux
	//handlers handlers.Handlers

	log *logrus.Entry
	cfg config.Config
}

func NewAPI(cfg config.Config, log *logrus.Logger, app *app.App) Api {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	return Api{
		//handlers: hands,
		//router:   router,
		server: server,
		log:    logger,
		cfg:    cfg,
	}
}

func (a *Api) Run(ctx context.Context, log *logrus.Logger) {

}
