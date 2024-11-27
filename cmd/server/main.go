package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	h "gowebsockets/internal/historian"
	"gowebsockets/internal/wsserver"
)

const (
	addr        = "192.168.1.7:8080"
	maxMessages = 10
	ttl         = 10 * time.Second
)

func main() {
	historian := h.NewHistorian(maxMessages, ttl)
	wsSrv := wsserver.NewWsServer(addr, historian)
	log.Info("Started ws server")
	if err := wsSrv.Start(); err != nil {
		log.Errorf("Error with ws server: %v", err)
	}
	log.Error(wsSrv.Stop())
}
