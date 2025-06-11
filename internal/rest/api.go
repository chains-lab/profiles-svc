package rest

import (
	"context"
	"net/http"

	"github.com/chains-lab/gatekit/mdlv"
	"github.com/chains-lab/gatekit/roles"
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

func NewRest(cfg config.Config, log *logrus.Logger, app *app.App) Api {
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
	auth := mdlv.AuthMdl(a.cfg.JWT.AccessToken.SecretKey, "todo")
	admin := mdlv.AccessGrant(a.cfg.JWT.AccessToken.SecretKey, "todo", roles.Admin, roles.SuperUser)
	verified := mdlv.AccessGrant(a.cfg.JWT.AccessToken.SecretKey, "todo", roles.Admin)
}
