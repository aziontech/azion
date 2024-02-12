package metric

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	analytics "github.com/segmentio/analytics-go/v3"
	"go.uber.org/zap"
)

var SegmentKey string = "YOUR_WRITE_KEY"

func location() string {
	dir, err := config.Dir()
	if err != nil {
		logger.Debug("Failed get path file metric.json", zap.Error(err))
	}
	const metricsFilename = "metrics.json"
	return filepath.Join(dir, metricsFilename)
}

func readLocalMetrics() map[string]int {
	file, err := os.OpenFile(location(), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}
	defer file.Close()

	var data map[string]int
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil && err != io.EOF {
		return nil
	}

	if data == nil {
		data = make(map[string]int)
	}

	return data
}

func Send(settings *token.Settings) {
	client := analytics.New(SegmentKey)
	fmt.Println("here", SegmentKey)
	defer client.Close()

	metrics := readLocalMetrics()

	for event, times := range metrics {
		fmt.Println("só garantindo")

		client.Enqueue(analytics.Track{
			UserId: settings.ClientId,
			Event:  event,
			Properties: analytics.NewProperties().
				Set("email", settings.Email).
				Set("times executed", times),
		})
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
