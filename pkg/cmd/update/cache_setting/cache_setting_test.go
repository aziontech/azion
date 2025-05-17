package cachesetting

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

}
