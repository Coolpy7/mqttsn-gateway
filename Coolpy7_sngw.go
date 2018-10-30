package main

import (
	"errors"
	"flag"
	"gateway"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		confFile  = flag.String("c", "", "config file path")
		topicFile = flag.String("t", "", "predefined topic file path")
	)
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if dir == "/" {
		dir += "data"
	} else {
		dir += "/data"
	}
	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	if *confFile == "" {
		*confFile = dir + "/mqttsn.yml"
	}

	// parse config
	config, err := mqttsn.ParseConfig(*confFile)
	if err != nil {
		panic(err)
	}

	// parse topic file
	if *topicFile != "" {
		err = mqttsn.InitPredefinedTopic(*topicFile)
		if err != nil {
			panic(err)
		}
	}

	// initialize logger
	err = mqttsn.InitLogger(dir + "/" + config.LogFilePath)
	if err != nil {
		panic(err)
	}

	// create signal chan
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// create Gateway
	var gateway mqttsn.Gateway
	if config.IsAggregate {
		gateway = mqttsn.NewAggregatingGateway(config, signalChan)
	} else {
		gateway = mqttsn.NewTransparentGateway(config, signalChan)
	}

	// start server
	err = gateway.StartUp()
	if err != nil {
		log.Println(errors.New("ERROR : failed to StartUp gateway"))
	}
	return
}
