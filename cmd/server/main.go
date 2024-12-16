package main

import (
	log "github.com/sirupsen/logrus"
	h "gowebsockets/internal/historian"
	"gowebsockets/internal/wsserver"
	"time"
)

const (
	addr      = "192.168.1.7:8080"
	maxMsgNum = 3
	ttl       = time.Second * 20
)

func main() {
	historian := h.NewHistorian(maxMsgNum, ttl)
	go historian.Start()
	wsSrv := wsserver.NewWsServer(addr, historian)
	log.Info("Started ws server")
	if err := wsSrv.Start(); err != nil {
		log.Errorf("Error with ws server: %v", err)
	}
	log.Error(wsSrv.Stop())
}
