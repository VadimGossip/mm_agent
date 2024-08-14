package main

import (
	"context"
	"github.com/VadimGossip/mm_agent/internal/app"
	"github.com/sirupsen/logrus"
	"time"
)

var configDir = "config"
var appName = "MM Agent"

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx, appName, configDir, time.Now())
	if err != nil {
		logrus.Fatalf("failed to init app[%s]: %s", appName, err)
	}

	a.Run(ctx)
}
