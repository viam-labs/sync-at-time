package timesyncsensor

import (
	"context"
	"errors"
	"fmt"
	"time"

	sensor "go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

var (
	Timesyncsensor   = resource.NewModel("naomi", "sync-at-time", "timesyncsensor")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(sensor.API, Timesyncsensor,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newTimeSyncer,
		},
	)
}

type Config struct {
    Start   string `json:"start"`
    End     string `json:"end"`
    Zone    string `json:"zone"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
    if cfg.Start == "" {
        return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "start")
    }

    if cfg.End == "" {
        return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "end")
    }

    if cfg.Zone == "" {
        return nil, nil, utils.NewConfigValidationFieldRequiredError(path, "zone")
    }

    return []string{}, []string{}, nil
}

type timeSyncer struct {
	resource.AlwaysRebuild

	name resource.Name
    start      string
    end        string
    zone       string

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
}

func newTimeSyncer(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		logger.Warn("Error configuring module with ", err)
		return nil, err
	}

	return NewTimeSyncSensor(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewTimeSyncSensor(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (sensor.Sensor, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &timeSyncer{
		name:       name,
		start:      conf.Start,
		end:        conf.End,
		zone:       conf.Zone,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	return s, nil
}

func (s *timeSyncer) Name() resource.Name {
	return s.name
}

func (s *timeSyncer) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
    currentTime := time.Now()
    var hStart, mStart, sStart, hEnd, mEnd, sEnd int
    n, err := fmt.Sscanf(s.start, "%d:%d:%d", &hStart, &mStart, &sStart)

    if err != nil || n != 3 {
        s.logger.Error("Start time is not in the format HH:MM:SS.")
        return nil, err
    }
    m, err := fmt.Sscanf(s.end, "%d:%d:%d", &hEnd, &mEnd, &sEnd)
    if err != nil || m != 3 {
        s.logger.Error("End time is not in the format HH:MM:SS.")
        return nil, err
    }

    zone, err := time.LoadLocation(s.zone)
    if err != nil {
        s.logger.Error("Time zone cannot be loaded: ", s.zone)
    }

    startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
        hStart, mStart, sStart, 0, zone)
    endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
        hEnd, mEnd, sEnd, 0, zone)

    readings := map[string]interface{}{"should_sync": false}
    readings["time"] = currentTime.String()
    // If it is between the start and end time, sync.
    if currentTime.After(startTime) && currentTime.Before(endTime) {
        s.logger.Debug("Syncing")
        readings["should_sync"] = true
        return readings, nil
    }

    // Otherwise, do not sync.
    s.logger.Debug("Not syncing. Current time not in sync window: " + currentTime.String())
    return readings, nil
}

func (s *timeSyncer) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *timeSyncer) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
