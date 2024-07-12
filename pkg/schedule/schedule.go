package schedule

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

type Schedule struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"` // schedule creation time
	Kind string    `json:"kind"`
}

var (
	TriggerDelete func(f *cmdutil.Factory, name string) error = TriggerDeleteBucket
)

func NewSchedule(name string, kind string) error {
	schedule := Schedule{
		Name: name,
		Time: time.Now(),
		Kind: kind,
	}

	schedules, err := readFileSchedule()
	if err != nil {
		logger.Debug("Error while reading the schedule", zap.Error(err))
		return err
	}
	schedules = append(schedules, schedule)

	err = createFileSchedule(schedules)
	if err != nil {
		logger.Debug("Scheduling error", zap.Error(err))
		return err
	}
	return nil
}

func createFileSchedule(shedule []Schedule) error {
	b, err := json.MarshalIndent(shedule, "	", " ")
	if err != nil {
		return err
	}
	configPath, err := config.Dir()
	if err != nil {
		return err
	}
	path := filepath.Join(configPath.Dir, configPath.Schedule)
	return os.WriteFile(path, b, os.FileMode(os.O_CREATE))
}

func readFileSchedule() ([]Schedule, error) {
	configPath, err := config.Dir()
	if err != nil {
		return nil, err
	}

	schedules := []Schedule{}

	path := filepath.Join(configPath.Dir, configPath.Schedule)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, err := json.Marshal(&schedules)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(path, data, 0666)
		if err != nil {
			return nil, err
		}

		return schedules, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(file) == 0 {
		return schedules, nil
	}

	err = json.Unmarshal(file, &schedules)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func ExecSchedules(factory *cmdutil.Factory) {
	logger.Debug("Exec Schedules")
	schedules, err := readFileSchedule()
	if err != nil {
		logger.Debug("Error while reading the schedule", zap.Error(err))
		return
	}

	scheds := []Schedule{}
	for _, s := range schedules {
		if CheckIf24HoursPassed(s.Time) {
			if s.Kind == DELETE_BUCKET {
				if err := TriggerDelete(factory, s.Name); err != nil {
					logger.Debug("Event execution error", zap.Error(err))
					scheds = append(scheds, s)
				}
			}
		}
	}

	if err := createFileSchedule(scheds); err != nil {
		logger.Debug("Scheduling error", zap.Error(err))
	}
}

// CheckIf24HoursPassed Checks if the current time is before 24 hours after the time 's'.
func CheckIf24HoursPassed(passed time.Time) bool {
	now := time.Now()
	diff := now.Sub(passed)
	return diff >= 24*time.Hour
}
