package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/deep-compute/docker-file-log-driver/driver"
	"github.com/deep-compute/docker-file-log-driver/http"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func main() {
	levelVal := os.Getenv("LOG_LEVEL")
	if levelVal == "" {
		levelVal = "info"
	}
	if level, exists := logLevels[levelVal]; exists {
		logrus.SetLevel(level)
	} else {
		fmt.Fprintln(os.Stderr, "invalid log level: ", levelVal)
		os.Exit(1)
	}

	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	http.Handlers(&h, driver.NewDriver())
	if err := h.ServeUnix("log", 0); err != nil {
		panic(err)
	}
}
