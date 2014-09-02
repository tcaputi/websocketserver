package websocketserver

import(
	//"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
var h = hub{
	connections: make(map[string]*connection),
	inbox:   make(chan *message),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
}

func RegisterReceiveHandler(cb func(string, string) string){
	h.receiveHandler = cb;
}

func Send(connid string, data []byte){
	h.connections[connid].outbox <- data
}

func Initialize(route string){
	http.HandleFunc(route, wsHandler)
	go h.run()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {	
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
		fmt.Println(ws, err)
        return
    }
    c := &connection{id: "1"/*uuid.New()*/, ws: ws, outbox: make(chan []byte, 256)}
    h.register <- c
    defer func() { h.unregister <- c }()
    go c.writer()
    c.reader()
}