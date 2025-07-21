package infra

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger *logrus.Logger

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}

	// Setup log directory
	logDir := "log/"
	if exists, err := dirExists(logDir); err != nil {
		panic(err)
	} else if !exists {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// Setup log file rotation
	writer, err := rotatelogs.New(
		logDir+viper.GetString("APP.NAME")+"-%Y%m%d.log",
		// rotatelogs.WithMaxAge(7*24*time.Hour),     // Keep logs for 7 days
		rotatelogs.WithRotationCount(7),           // Keep last 7 files
		rotatelogs.WithRotationTime(24*time.Hour), // Rotate daily
	)
	if err != nil {
		panic(err)
	}

	// Instantiate logger
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// JSON Formatter with pretty indentation (human readable)
	jsonFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true, // This makes it "prettier"
	}

	// Hook untuk tulis ke file (JSON log)
	logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.DebugLevel: writer,
			logrus.WarnLevel:  writer,
		},
		jsonFormatter,
	))

	// Stdout Formatter
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(jsonFormatter)

	return logger
}
