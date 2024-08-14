package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/VadimGossip/platform_common/pkg/closer"
	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/mm_agent/internal/config"
	"github.com/VadimGossip/mm_agent/internal/model"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type App struct {
	serviceProvider *serviceProvider
	name            string
	configDir       string
	appStartedAt    time.Time
	cfg             *model.Config
}

func NewApp(ctx context.Context, name, configDir string, appStartedAt time.Time) (*App, error) {
	a := &App{
		name:         name,
		configDir:    configDir,
		appStartedAt: appStartedAt,
	}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.Init(a.configDir)
	if err != nil {
		return fmt.Errorf("[%s] config initialization error: %s", a.name, err)
	}
	a.cfg = cfg
	logrus.Infof("[%s] got config: [%+v]", a.name, *a.cfg)
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.cfg)
	return nil
}

func (a *App) Run(ctx context.Context) {
	defer func() {
		closer.CloseAll()
		closer.Wait()
		logrus.Infof("[%s] stopped", a.name)
	}()
	logrus.Infof("[%s] started", a.name)
	a.serviceProvider.MsgService(ctx).RunScanner(ctx)
}
