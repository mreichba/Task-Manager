package logger

import (
	"os"
	"strings"

	"github.com/mreichba/task-manager-backend/config"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetOutput(os.Stdout)

	env := strings.ToLower(config.AppConfig.Environment)

	switch env {
	case "production":
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.JSONFormatter{})
	case "test":
		Log.SetLevel(logrus.WarnLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			DisableQuote:  true,
			FullTimestamp: true,
		})
	default: // development, staging, etc.
		Log.SetLevel(logrus.DebugLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}
}

func Info(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Info(msg)
}

func Warn(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Warn(msg)
}

func Error(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Error(msg)
}

func Fatal(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Fatal(msg)
}

func Debug(msg string, fields logrus.Fields) {
	Log.WithFields(fields).Debug(msg)
}
