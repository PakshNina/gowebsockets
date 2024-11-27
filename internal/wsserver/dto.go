package wsserver

type wsMessage struct {
	IPAddress string `json:"address"`
	Message   string `json:"message"`
	Time      string `json:"time"`
}
