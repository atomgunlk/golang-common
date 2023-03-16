package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/atomgunlk/golang-common/pkg/logger/hooks"
	"github.com/sirupsen/logrus"
)

// TimeStampFormat config time format
const TimeStampFormat = "2006-01-02 15:04:05.000"

// Fields sub coordinate with WithField method
type Fields map[string]interface{}

// Entry map key value for logger.WithFields
type Entry struct {
	*logrus.Entry
}

var (
	logger                          = logrus.New()
	mw                    io.Writer = nil
	logStdOut             io.Writer = nil
	logToStdout           bool      = true
	logPath               string    = "log"
	logFileNamePrefix     string    = "log"
	logFileNameDateSuffix bool      = true
	logFile               *os.File  = nil
	lastLogTime           time.Time
	lastLogFileOpenTime   time.Time
	logfileCloseTimeout   int    = 10 // in second
	appenv                string = os.Getenv("APP_ENV")
)

func init() {
	switch runtime.GOOS {
	case "darwin":
		logStdOut = os.Stdout
	case "windows":
		logStdOut = os.Stdout
		// logStdOut = ansicolor.NewAnsiColorWriter(os.Stdout)
	default:
		logStdOut = os.Stdout
	}

	if appenv != "production" || appenv == "" {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: false,
			DisableColors:             false,
			TimestampFormat:           TimeStampFormat,
			FullTimestamp:             true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "log_level",
			},
		})

		level := getEnv("LOG_LEVEL", "debug")
		SetLevel(level)
		// hooks.UseTimestamp(logger)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "log_level",
			},
			PrettyPrint: false,
		})

		level := getEnv("LOG_LEVEL", "info")
		SetLevel(level)
		// logger.SetOutput(os.Stdout)
		hooks.UseTimestamp(logger)
	}
}
func GetLogger() *logrus.Logger {
	return logger
}

// SetOutputToFile
// out to file
func SetOutputToFile(path, filenamePrefix string, useDateSuffix bool) {
	logToStdout = false

	logPath = path
	logFileNamePrefix = filenamePrefix
	logFileNameDateSuffix = useDateSuffix

	// upload path
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		fmt.Printf("Logger Create log Path %s\r\n", err.Error())
	}
	// go routine close file timer
	if !isLogfileOpen() {
		logFileOpen()
		// go logFileCloseTimer()
	}
}

// SetOutputToStdOut
// out to os.StdOut
func SetOutputToStdOut() {
	logToStdout = true
	logger.SetOutput(os.Stdout)
}

// SetOutputToMW , Multi Writer
// out to os.StdOut and File
func SetOutputToMW(path, filenamePrefix string, useDateSuffix bool) {
	logToStdout = true

	logPath = path
	logFileNamePrefix = filenamePrefix
	logFileNameDateSuffix = useDateSuffix

	// upload path
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		fmt.Printf("Logger Create log Path %s\r\n", err.Error())
	}
	// go routine close file timer
	if !isLogfileOpen() {
		logFileOpen()
		// go logFileCloseTimer()
	}
}

// AddFields appends fields to the log entry using custom logger.Fields
func AddFields(e *Entry, f Fields) *Entry {
	return &Entry{e.WithFields(logrus.Fields(f))}
}

// AddHook with third-party integration
func AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// SetLevel of logging level
func SetLevel(level string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Level = logrus.InfoLevel
	} else {
		logger.Level = l
	}
}

// DisableColor for logging console
func DisableColor() {
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   true,
		TimestampFormat: TimeStampFormat,
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Infof(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Printf(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
// it do defer os.Exit(1) after print log
func Fatalf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Fatalf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
// it do panic() after print log
func Panicf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Panicf(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Debug(args...)
}

// Info - print information to stdout
// You can print output info you need.
//
//	func ExampleExamples_output() {
//	    logger.Info("Hello")
//	    // Output: Hello
//	}
func Info(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Info(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Print(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Fatal(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Panic(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Debugln(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Infoln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Println(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Warning(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Errorln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Fatalln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": appenv,
	})

	logger.Panicln(args...)
}

// WithFields log with your specific field
func WithFields(f Fields) *Entry {
	f["environment"] = appenv

	return &Entry{logger.WithFields(logrus.Fields(f))}
}

// WithError log with error
func WithError(err error) *Entry {
	return &Entry{logger.WithError(err)}
}

// StandardLogger for gin gonic integration with logger
func StandardLogger() *logrus.Logger {
	return logger
}
