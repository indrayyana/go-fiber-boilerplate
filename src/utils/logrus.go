package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	logrus.TextFormatter
}

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Set logger to use the custom text formatter
	Log.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: "15:04:05.000",
			FullTimestamp:   true,
			ForceColors:     true,
		},
	})

	Log.SetOutput(os.Stdout)
}
