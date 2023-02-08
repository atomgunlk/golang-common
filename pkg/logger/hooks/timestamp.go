package hooks

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	logFieldTimestamp = "timestamp"
)

type timestampHook struct{}

func (t timestampHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t timestampHook) Fire(entry *logrus.Entry) error {
	entry.Data[logFieldTimestamp] = time.Now().UnixNano()

	return nil
}

// UseTimestamp adds timestamp to logs
func UseTimestamp(logger *logrus.Logger) {
	hook := timestampHook{}
	logger.Hooks.Add(hook)
}
