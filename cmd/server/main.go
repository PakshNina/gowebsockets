package main

import (
	log "github.com/sirupsen/logrus"
	"gowebsockets/internal/wsserver"
)

const (
	addr = "192.168.1.7:8080"
)

func main() {
	wsSrv := wsserver.NewWsServer(addr)
	log.Info("Started ws server")
	if err := wsSrv.Start(); err != nil {
		log.Errorf("Error with ws server: %v", err)
	}
	log.Error(wsSrv.Stop())
}
