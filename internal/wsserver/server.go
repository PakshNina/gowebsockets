package wsserver

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	templateDir = "./web/templates/html"
	staticDir   = "./web/static/"
)

type WSServer interface {
	Start() error
	Stop() error
}

type historian interface {
	SaveMessage(msg interface{})
	GetAllMessages() []interface{}
}

type wsSrv struct {
	mux       *http.ServeMux
	srv       *http.Server
	wsUpg     *websocket.Upgrader
	wsClients map[*websocket.Conn]struct{}
	mutex     *sync.RWMutex
	broadcast chan *wsMessage
	historian historian
}

func NewWsServer(addr string, h historian) WSServer {
	m := http.NewServeMux()
	return &wsSrv{
		mux: m,
		srv: &http.Server{
			Addr:    addr,
			Handler: m,
		},
		wsUpg:     &websocket.Upgrader{},
		wsClients: map[*websocket.Conn]struct{}{},
		mutex:     &sync.RWMutex{},
		broadcast: make(chan *wsMessage),
		historian: h,
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	ws.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	go ws.writeToClientsBroadcast()
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) Stop() error {
	log.Info("Before", ws.wsClients)
	close(ws.broadcast)
	ws.mutex.Lock()
	for conn := range ws.wsClients {
		if err := conn.Close(); err != nil {
			log.Errorf("Error with closing: %v", err)
		}
		delete(ws.wsClients, conn)
	}
	ws.mutex.Unlock()
	log.Info("After close", ws.wsClients)
	return ws.srv.Shutdown(context.Background())
}

func (ws *wsSrv) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error with websocket connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Infof("Client with address %s connected", conn.RemoteAddr().String())
	ws.mutex.Lock()
	ws.wsClients[conn] = struct{}{}
	ws.mutex.Unlock()
	ws.sendHistory(conn)
	go ws.readFromClient(conn)
}

func (ws *wsSrv) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(wsMessage)
		if err := conn.ReadJSON(msg); err != nil {
			wsErr, ok := err.(*websocket.CloseError)
			if !ok || wsErr.Code != websocket.CloseGoingAway {
				log.Errorf("Error with reading from WebSocket: %v", err)
			} else {
				log.Infof("Client %s disconnected", conn.RemoteAddr().String())
			}
			break
		}
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Errorf("Error with address split: %v", err)
		}
		msg.IPAddress = host
		msg.Time = time.Now().Format("15:04")
		go ws.historian.SaveMessage(msg)
		ws.broadcast <- msg
	}
	ws.mutex.Lock()
	delete(ws.wsClients, conn)
	ws.mutex.Unlock()
}

func (ws *wsSrv) writeToClientsBroadcast() {
	for msg := range ws.broadcast {
		ws.mutex.RLock()
		for client := range ws.wsClients {
			func() {
				if err := client.WriteJSON(msg); err != nil {
					log.Errorf("Error with writing message: %v", err)
				}
			}()
		}
		ws.mutex.RUnlock()
	}
}

func (ws *wsSrv) sendHistory(conn *websocket.Conn) {
	for _, msg := range ws.historian.GetAllMessages() {
		if err := conn.WriteJSON(msg); err != nil {
			log.Errorf("Error with sending history: %v", err)
		}
	}
}
