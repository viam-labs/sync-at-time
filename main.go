// Package main is a module which serves the timesync sensor model.
package main

import (
	"context"

	"go.viam.com/utils"

	timesyncsensor "github.com/viam-labs/sync-at-time/timesyncsensor"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("sync_at_time"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	myMod, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = myMod.AddModelFromRegistry(ctx, sensor.API, timesyncsensor.Model)
	if err != nil {
		return err
	}

	err = myMod.Start(ctx)
	defer myMod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}