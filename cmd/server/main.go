package main

import (
	"flag"
	v1 "github.com/StellrisJAY/workflow-ai/internal/api/v1"
	"github.com/StellrisJAY/workflow-ai/internal/config"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to config file")
)

func main() {
	flag.Parse()
	conf, err := config.ParseConfig(*configFile)
	if err != nil {
		panic(err)
	}
	router := v1.NewRouter(conf)
	if err := router.Init(); err != nil {
		panic(err)
	}
	if err = router.Start(); err != nil {
		panic(err)
	}
}
