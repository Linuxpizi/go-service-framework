package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	logFile := "stderr"

	err := Init(DebugLevel, logFile)
	if err != nil {
		t.Errorf("Fail to init logger with debug level: %v", err)
	}

	sugar := Sugar()
	if sugar == nil {
		t.Errorf("Sugar is not initialized")
	}

	if sugar != nil {
		sugar.Info("test for logger infow",
			zap.String("file", logFile),
		)
	}

	sugar.Info("test for logger infof",
		zap.String("file", logFile),
	)

	f := func() int {
		f, err := os.Open("xx")
		if err != nil {
			sugar.Error("open file",
				zap.String("filename", "xx"),
				zap.Error(err),
			)
		}
		defer f.Close()
		return 0
	}
	sugar.Error("test for logger infof",
		zap.String("file", logFile),
		zap.Int("res", f()),
	)
}
