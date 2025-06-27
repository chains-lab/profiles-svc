package app

import (
	"github.com/chains-lab/elector-cab-svc/internal/utils/config"
	"github.com/sirupsen/logrus"
)

type App struct {
}

func NewApp(cfg config.Config, log *logrus.Logger) (App, error) {
	return App{}, nil
}
