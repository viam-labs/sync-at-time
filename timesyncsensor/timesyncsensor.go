// Package timesyncsensor implements a datasync manager.
package timesyncsensor

import (
    "context"
    "fmt"
    "errors"
    "time"

    "go.viam.com/rdk/components/sensor"
    "go.viam.com/rdk/logging"
    "go.viam.com/rdk/resource"
    "go.viam.com/rdk/services/datamanager"

    "go.viam.com/utils"
)

// Model is the full model definition.
var (
    Model            = resource.NewModel("naomi", "sync-at-time", "timesyncsensor")
    errUnimplemented = errors.New("unimplemented")
)

func init() {
    resource.RegisterComponent(sensor.API, Model,
        resource.Registration[sensor.Sensor, *Config]{
            Constructor: func(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
                sensorConfig, err := resource.NativeConfig[*Config](conf)
                if err != nil {
                    logger.Warnf("Error configuring module with %+s", err)
                    return nil, err
                }

                cancelCtx, cancelFunc := context.WithCancel(context.Background())

                v := &timeSyncer{
                    name:       conf.ResourceName(),
                    start:      sensorConfig.Start,
                    end:        sensorConfig.End,
                    zone:       sensorConfig.Zone,
                    logger:     logger,
                    cancelCtx:  cancelCtx,
                    cancelFunc: cancelFunc,
                }
                if err := v.Reconfigure(ctx, deps, conf); err != nil {
                    logger.Warnf("Error configuring module with %+s", err)
                    return nil, err
                }
                return v, nil
            },
        })
}

type Config struct {
    Start   string `json:"start"`
    End     string `json:"end"`
    Zone    string `json:"zone"`
}

// Validate validates the config and returns implicit dependencies.
func (cfg *Config) Validate(path string) ([]string, error) {
    if cfg.Start == "" {
        return nil, utils.NewConfigValidationFieldRequiredError(path, "start")
    }

    if cfg.End == "" {
        return nil, utils.NewConfigValidationFieldRequiredError(path, "end")
    }

    if cfg.Zone == "" {
        return nil, utils.NewConfigValidationFieldRequiredError(path, "zone")
    }

    return []string{}, nil
}


type timeSyncer struct {
    name       resource.Name
    start      string
    end        string
    zone       string

    cancelCtx  context.Context
    cancelFunc func()
    logger     logging.Logger
}

func (s *timeSyncer) Name() resource.Name {
    return s.name
}

func (s *timeSyncer) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
    sensorConfig, err := resource.NativeConfig[*Config](conf)
    if err != nil {
        s.logger.Warnf("Error reconfiguring module with %+s", err)
        return nil
    }

    s.start = sensorConfig.Start
    s.end = sensorConfig.End
	s.name = conf.ResourceName()
    s.logger.Info("Start time and end time have changed to %s and %s.", s.start, s.end)

    return nil
}

// DoCommand does nothing.
func (s *timeSyncer) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
    return nil, errUnimplemented
}

func (s *timeSyncer) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
    currentTime := time.Now()
    var hStart, mStart, sStart, hEnd, mEnd, sEnd int
    n, err := fmt.Sscanf(s.start, "%d:%d:%d", &hStart, &mStart, &sStart)
    if err != nil || n != 3 {
        s.logger.Error("Start time is not in the format HH:MM:SS.")
        return datamanager.CreateShouldSyncReading(false), err
    }
    m, err := fmt.Sscanf(s.end, "%d:%d:%d", &hEnd, &mEnd, &sEnd)
    if err != nil || m != 3 {
        s.logger.Error("End time is not in the format HH:MM:SS.")
        return datamanager.CreateShouldSyncReading(false), err
    }

    zone, err := time.LoadLocation(s.zone)
    if err != nil {
        s.logger.Error("Time zone %s cannot be loaded", s.zone)
    }

    startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
        hStart, mStart, sStart, 0, zone)
    endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
        hEnd, mEnd, sEnd, 0, zone)

    // If it is between the start and end time, sync.
    if currentTime.After(startTime) && currentTime.Before(endTime) {
        s.logger.Info("Syncing")
        return datamanager.CreateShouldSyncReading(true), nil
    }

    // Otherwise, do not sync.
    return datamanager.CreateShouldSyncReading(false), nil
}

// Close closes the underlying generic.
func (s *timeSyncer) Close(ctx context.Context) error {
    s.cancelFunc()
    return nil
}