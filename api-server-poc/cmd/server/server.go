package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	app "api-server-poc"
	"api-server-poc/config"
	"api-server-poc/logger"
	"api-server-poc/server"

	"github.com/sirupsen/logrus"
)

type closer func()

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	conf := config.LoadConfigFile()
	close := loggerInit(conf)
	defer close()

	app := app.NewDIContainer()

	app.Provide(
		config.LoadConfigFile,
	)

	app.Invoke(
		server.GRPCServer,
		server.APIServer,
	)

	app.Run()
}

func loggerInit(conf *config.Config) closer {
	var logOutput io.Writer
	var output closer = func() {}
	if conf.Log.Type == "file" {
		f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(fmt.Sprintf("error opening log file: %v", err))
		}
		logOutput = f
		output = func() {
			f.Close()
		}
	} else {
		logOutput = os.Stdout
	}
	level, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		panic(fmt.Sprintf("error log level : %v", err))
	}
	logger.Init(logOutput, level)
	return output
}
