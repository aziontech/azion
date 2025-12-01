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

type Factory struct {
	Schedule      Schedule
	Dir           func() config.DirPath
	Join          func(elem ...string) string
	Stat          func(name string) (os.FileInfo, error)
	IsNotExist    func(err error) bool
	Marshal       func(v any) ([]byte, error)
	WriteFile     func(name string, data []byte, perm os.FileMode) error
	ReadFile      func(name string) ([]byte, error)
	Unmarshal     func(data []byte, v any) error
	MarshalIndent func(v any, prefix, indent string) ([]byte, error)
}

var factoryShedule Factory

func init() {
	InjectFactory(nil)
}

func InjectFactory(fact *Factory) *Factory {
	if fact != nil {
		factoryShedule = *fact
		return fact
	}

	f := &Factory{
		Dir:           config.Dir,
		Join:          filepath.Join,
		Stat:          os.Stat,
		IsNotExist:    os.IsNotExist,
		Marshal:       json.Marshal,
		WriteFile:     os.WriteFile,
		ReadFile:      os.ReadFile,
		Unmarshal:     json.Unmarshal,
		MarshalIndent: json.MarshalIndent,
	}
	factoryShedule = *f
	return f
}

type Schedule struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"` // schedule creation time
	Kind string    `json:"kind"`
}

func NewSchedule(fact *Factory, cmdFactory *cmdutil.Factory, name string, kind string) error {
	factory := InjectFactory(fact)
	factory.Schedule.Name = name
	factory.Schedule.Time = time.Now()
	factory.Schedule.Kind = kind

	var activeProfile string
	if cmdFactory != nil {
		activeProfile = cmdFactory.GetActiveProfile()
	}

	schedules, err := factory.readFileScheduleForProfile(activeProfile)
	if err != nil {
		logger.Debug("Error while reading the schedule", zap.Error(err))
		return err
	}

	schedules = append(schedules, factory.Schedule)

	err = factory.createFileScheduleForProfile(schedules, activeProfile)
	if err != nil {
		logger.Debug("Scheduling error", zap.Error(err))
		return err
	}
	return nil
}

func (s *Factory) createFileScheduleForProfile(schedules []Schedule, profile string) error {
	b, err := s.MarshalIndent(schedules, "  ", " ")
	if err != nil {
		return err
	}
	configPath := s.Dir()
	if profile != "" {
		configPath.Dir = filepath.Join(configPath.Dir, profile)
	}
	path := s.Join(configPath.Dir, configPath.Schedule)
	return s.WriteFile(path, b, 0666)
}

func (s *Factory) readFileScheduleForProfile(profile string) ([]Schedule, error) {
	configPath := config.Dir()
	if profile != "" {
		configPath.Dir = filepath.Join(configPath.Dir, profile)
	}
	schedules := []Schedule{}

	path := s.Join(configPath.Dir, configPath.Schedule)

	// Checks if the file exists in the given path
	if _, err := s.Stat(path); s.IsNotExist(err) {
		data, err := s.Marshal(&schedules)
		if err != nil {
			return nil, err
		}

		if err := s.WriteFile(path, data, 0666); err != nil {
			return nil, err
		}
		return schedules, nil
	}

	file, err := s.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = s.Unmarshal(file, &schedules)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func ExecSchedules(factory *cmdutil.Factory) {
	logger.Debug("Exec Schedules")
	activeProfile := factory.GetActiveProfile()
	schedules, err := factoryShedule.readFileScheduleForProfile(activeProfile)
	if err != nil {
		logger.Debug("Error while reading the schedule", zap.Error(err))
		return
	}

	scheds := []Schedule{}
	for _, s := range schedules {
		if CheckIf24HoursPassed(s.Time) {
			if s.Kind == DELETE_BUCKET {
				if err := TriggerDeleteBucket(factory, s.Name); err != nil {
					logger.Debug("Event execution error", zap.Error(err))
					scheds = append(scheds, s)
				}
			}
		}
	}

	if err := factoryShedule.createFileScheduleForProfile(scheds, activeProfile); err != nil {
		logger.Debug("Scheduling error", zap.Error(err))
	}
}

// CheckIf24HoursPassed Checks if the current time is before 24 hours after the time 's'.
func CheckIf24HoursPassed(passed time.Time) bool {
	now := time.Now()
	diff := now.Sub(passed)
	return diff >= 24*time.Hour
}
