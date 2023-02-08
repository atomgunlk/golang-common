package logger

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/atomgunlk/golang-common/pkg/logger/hooks"
)

// TimeStampFormat config time format
const TimeStampFormat = "2006-01-02 15:04:05.000"

// Fields sub coordinate with WithField method
type Fields map[string]interface{}

// Entry map key value for logger.WithFields
type Entry struct {
	*logrus.Entry
}

var logger = logrus.New()

func init() {
	if os.Getenv("APP_ENV") != "production" || os.Getenv("APP_ENV") == "" {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			DisableColors:             false,
			TimestampFormat:           TimeStampFormat,
			FullTimestamp:             true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "log_level",
			},
		})

		level := getEnv("LOG_LEVEL", "debug")
		SetLevel(level)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "log_level",
			},
			PrettyPrint: true,
		})

		level := getEnv("LOG_LEVEL", "info")
		SetLevel(level)
		logger.SetOutput(os.Stdout)
		hooks.UseTimestamp(logger)
	}
}
func GetLogger() *logrus.Logger {
	return logger
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
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Infof(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Printf(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
// it do defer os.Exit(1) after print log
func Fatalf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Fatalf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
// it do panic() after print log
func Panicf(format string, args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Panicf(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
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
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Info(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Print(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Fatal(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Panic(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Debugln(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Infoln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Println(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Warning(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Errorln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Fatalln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logger := logger.WithFields(logrus.Fields{
		"environment": os.Getenv("APP_ENV"),
	})

	logger.Panicln(args...)
}

// WithFields log with your specific field
func WithFields(f Fields) *Entry {
	f["environment"] = os.Getenv("APP_ENV")

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
