package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	AutoLogFile   = "auto"
	logDateFormat = "2006-01-02T15-04-05"
)

// FileHook to send logs to the trace file regardless of CLI level.
type FileHook struct {
	// Logger is a reference to the internal Logger that this utilizes.
	Logger *logrus.Logger

	// File is a reference to the file being written to.
	File *os.File

	// FilePath is the full path used when creating the file.
	FilePath string

	// Closed reports whether the hook has been closed. Once closed
	// it can't be written to again.
	Closed bool
}

func NewFileHook(path string) (*FileHook, error) {
	if path == AutoLogFile {
		path = filepath.Join(os.Getenv("HOME"), ".crashd", getLogNameFromTime(time.Now()))
	}
	file, err := os.Create(path)
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	logger.SetOutput(file)

	logrus.Infof("Detailed logs being written to: %v", path)

	return &FileHook{Logger: logger, File: file, FilePath: path}, err
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	if hook.Closed {
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

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// CloseFileHooks will close each file being used for each FileHook attached to the logger.
// If the logger passed is nil, will reference the logrus.StandardLogger().
func CloseFileHooks(l *logrus.Logger) error {
	// All the hooks we utilize are just tied to the standard logger.
	logrus.Debugln("Closing log file; future log calls will not be persisted.")
	if l == nil {
		l = logrus.StandardLogger()
	}

	for _, fh := range GetFileHooks(l) {
		fh.File.Close()
		fh.Closed = true
	}
	return nil
}

// GetFirstFileHook is a convenience method to take an object and returns the first
// FileHook attached to it. Accepts an interface{} since the logger objects may be put
// into context or thread objects. The obj should be a *logrus.Logger object.
func GetFirstFileHook(obj interface{}) *FileHook {
	fhs := GetFileHooks(obj)
	if len(fhs) > 0 {
		return fhs[0]
	}
	return nil
}

// GetFileHooks is a convenience method to take an object and returns the
// FileHooks attached to it. Accepts an interface{} since the logger objects may be put
// into context or thread objects. The obj should be a *logrus.Logger object.
func GetFileHooks(obj interface{}) []*FileHook {
	l, ok := obj.(*logrus.Logger)
	if !ok {
		return nil
	}
	result := []*FileHook{}
	for _, hooks := range l.Hooks {
		for _, hook := range hooks {
			switch fh := hook.(type) {
			case *FileHook:
				result = append(result, fh)
			}
		}
	}
	return result
}

func getLogNameFromTime(t time.Time) string {
	return fmt.Sprintf("crashd_%v.log", t.Format(logDateFormat))
}
