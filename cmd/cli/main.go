package main

import (
	"context"
	"sync-at-time-2"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	sensor "go.viam.com/rdk/components/sensor"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	ctx := context.Background()
	logger := logging.NewLogger("cli")

	deps := resource.Dependencies{}
	// can load these from a remote machine if you need

	cfg := syncattime2.Config{}

	thing, err := syncattime2.NewTimesyncsensor(ctx, deps, sensor.Named("foo"), &cfg, logger)
	if err != nil {
		return err
	}
	defer thing.Close(ctx)

	return nil
}
