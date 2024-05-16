package schedule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

type Schedule struct {
	Name   string           `json:"name"`
	Time   time.Time        `json:"time"` // schedule creation time
	Object SheduleInterface `json:"object"`
}

type SheduleInterface interface {
	TriggerEvent() error
}

func NewSchedule(name string, object SheduleInterface) error {
	schedule := Schedule{
		Name:   name,
		Time:   time.Now(),
		Object: object,
	}

	schedules, err := readFileShedule()
	if err != nil {
		logger.Debug("read shedule error", zap.Error(err))
		return err
	}
	schedules = append(schedules, schedule)

	err = createFileShedule(schedules)
	if err != nil {
		logger.Debug("scheduling error", zap.Error(err))
		return err
	}
	return nil
}

func createFileShedule(shedule []Schedule) error {
	fmt.Println(">> here: ", shedule)
	b, err := json.MarshalIndent(shedule, "	", " ")
	if err != nil {
		fmt.Println(">>> err:", err.Error())
		return err
	}
	fmt.Println(">> b: ", b)
	configPath, err := config.Dir()
	if err != nil {
		return err
	}
	path := filepath.Join(configPath.Dir, configPath.Schedule)
	fmt.Println(">> path: ", path)
	return os.WriteFile(path, b, os.FileMode(os.O_CREATE))
}

func readFileShedule() ([]Schedule, error) {
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

func ExecSchedules() {
	schedules, err := readFileShedule()
	if err != nil {
		logger.Debug("read shedule error", zap.Error(err))
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(schedules))
	errChan := make(chan error, len(schedules))

	for _, s := range schedules {
		go func(s Schedule) {
			defer wg.Done()

			// Checks if the current time is before 24 hours after the time 's'.
			if time.Now().Before(s.Time.Add(24*time.Hour)) {
				if err := s.Object.TriggerEvent(); err != nil {
					errChan <- err
				}
			}
		}(s)
	}
	wg.Wait()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	go func() {
		for err := range errChan {
			logger.Debug("event execution error", zap.Error(err))
		}
	}()
}
