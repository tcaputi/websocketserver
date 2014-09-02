package websocketserver

import(
	"github.com/gorilla/websocket"
	"fmt"
)

type hub struct {
    connections map[string]*connection
    inbox chan *message
    register chan *connection
    unregister chan *connection
	receiveHandler func(string, string) string
}

type connection struct {
	id string
    ws *websocket.Conn
    outbox chan []byte
}

type message struct {
	connid string
	data []byte
}

func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
			fmt.Println("registered")
            h.connections[c.id] = c 
        case c := <-h.unregister:
			fmt.Println("unregistered")
            if _, ok := h.connections[c.id]; ok {
                delete(h.connections, c.id)
                close(c.outbox)
            }
        case m := <-h.inbox:
            h.connections[m.connid].outbox <- []byte(h.receiveHandler(m.connid, string(m.data)))
        }
    }
}

func (c *connection) reader() {
    for {
        _, data, err := c.ws.ReadMessage()
        if err != nil {
            break
        }
        h.inbox <- &message{c.id, data}
    }
    c.ws.Close()
}

func (c *connection) writer() {
    for data := range c.outbox {
        err := c.ws.WriteMessage(websocket.TextMessage, data)
        if err != nil {
            break
        }
    }
    c.ws.Close()
}