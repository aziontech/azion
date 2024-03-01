package metric

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	analytics "github.com/segmentio/analytics-go/v3"
	"go.uber.org/zap"
)

const SEGMENT_KEY = "Irg63QfdvWpoANAVeCBEwfxXBKvoSSzt"

func location() string {
	dir, err := config.Dir()
	if err != nil {
		logger.Debug("Failed get path file metric.json", zap.Error(err))
	}
	return filepath.Join(dir.Dir, dir.Metrics)
}

func readLocalMetrics() map[string]command {
	file, err := os.OpenFile(location(), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}
	defer file.Close()

	var data map[string]command
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil && err != io.EOF {
		return nil
	}

	if data == nil {
		data = make(map[string]command)
	}

	return data
}

func Send(settings *token.Settings) {
	client := analytics.New(SEGMENT_KEY)
	defer client.Close()

	metrics := readLocalMetrics()

	os := runtime.GOOS
	arch := runtime.GOARCH

	for event, cmd := range metrics {
		err := client.Enqueue(analytics.Track{
			UserId: settings.ClientId,
			Event:  fmt.Sprintf("cli_%s", event),
			Properties: analytics.NewProperties().
				Set("email", settings.Email).
				Set("cli version", cmd.CLIVersion).
				Set("vulcan version", cmd.VulcanVersion).
				Set("total successful", cmd.TotalSuccess).
				Set("total failed", cmd.TotalFailed).
				Set("total", cmd.TotalSuccess+cmd.TotalFailed).
				Set("shell", cmd.Shell).
				Set("execution time", cmd.ExecutionTime).
				Set("operating system", os).
				Set("architecture", arch).
				Set("client id", settings.ClientId),
		})
		if err != nil {
			logger.Debug("failed to send metrics", zap.Error(err))
			return
		}
	}

	clean()
}

// cleans metrics location and rewrites the file with empty content
func clean() {
	err := os.WriteFile(location(), []byte{}, 0666)
	if err != nil {
		return
	}
}
