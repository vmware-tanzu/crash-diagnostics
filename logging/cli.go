package logging

import (
	"io"

	"github.com/sirupsen/logrus"
)

// CLIHook creates returns a hook which will log at the specified level.
type CLIHook struct {
	Logger *logrus.Logger
	level  logrus.Level
}

func NewCLIHook(w io.Writer, level logrus.Level) *CLIHook {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(w)

	return &CLIHook{logger, level}
}

func (hook *CLIHook) Fire(entry *logrus.Entry) error {
	if !hook.Logger.IsLevelEnabled(entry.Level) {
		return nil
	}
	switch entry.Level {
	case logrus.PanicLevel:
		hook.Logger.Panic(entry.Message)
	case logrus.FatalLevel:
		hook.Logger.Fatal(entry.Message)
	case logrus.ErrorLevel:
		hook.Logger.Error(entry.Message)
	case logrus.WarnLevel:
		hook.Logger.Warning(entry.Message)
	case logrus.InfoLevel:
		hook.Logger.Info(entry.Message)
	case logrus.DebugLevel:
		hook.Logger.Debug(entry.Message)
	case logrus.TraceLevel:
		hook.Logger.Trace(entry.Message)
	default:
		hook.Logger.Info(entry.Message)
	}
	return nil
}

func (hook *CLIHook) Levels() []logrus.Level {
	output := []logrus.Level{}
	for _, v := range logrus.AllLevels {
		if hook.Logger.IsLevelEnabled(v) {
			output = append(output, v)
		}
	}
	return output
}
