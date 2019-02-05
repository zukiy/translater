package main

import (
	"flag"
	"os"
	"os/signal"
	"translator/api/v1"

	"translator/config"
	"translator/model"
	"translator/providers/storage"
	"translator/utils"
)

var (
	configFile  = flag.String("config", "config.yaml", "config file")
	interactive = flag.Bool("i", false, "interactive mode")
	httpPort    = flag.Int("httpPort", 8008, "HTTP server port")
)

func main() {
	flag.Parse()

	c := &config.Config{}

	err := utils.MustReadYAML(*configFile, &c)
	if err != nil {
		panic(err)
	}

	st := storage.Init(c)

	httpServ := v1.New(*httpPort, st, model.New(st))
	go httpServ.Serve()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	httpServ.Stop()

	println("exit")
}
