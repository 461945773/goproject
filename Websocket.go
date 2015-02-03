// Websocket.go
package main

import (
	"fmt"
	"net/http"
	// "text/template"
	// "encoding/json"
	"websocket"
)

/*Global*/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type dataGroup struct {
	DataType int
	Sender   string
	Message  string
}

const (
	LINE_SEGMENT = iota
	CHAT_MESSAGE
	GAME_LOGIC
)

// var indexTemp = template.Must(template.ParseFiles("index.html"))

/*Hub*/
type hub struct {
	connections map[*connection]bool
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

var Hub = hub{
	connections: make(map[*connection]bool),
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}

/*connection*/

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *connection) reader() {
	defer c.ws.Close()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		Hub.broadcast <- message
		fmt.Println(string(message))
	}
}

func (c *connection) writer() {
	defer c.ws.Close()
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

func serverWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	fmt.Println(r.Header.Get("Origin"), r.Host)
	if err != nil {
		fmt.Println("conn", err)
		return
	}
	c := &connection{
		send: make(chan []byte, 1024), ws: ws}
	Hub.register <- c
	defer func() {
		Hub.unregister <- c
	}()
	// data := dataGroup{
	// 	DataType: CHAT_MESSAGE,
	// 	Message:  "欢迎" + "fuck" + "来到你画我猜游戏大厅",
	// 	Sender:   "你画我猜",
	// }
	// b, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("json err", err)
	// }
	// Hub.broadcast <- b
	go c.writer()
	c.reader()
}

/*ServerHandler*/

// func serverIndex(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		http.Error(w, "Not found", 404)
// 		return
// 	}
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not allowed", 405)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	indexTemp.Execute(w, r.Host)
// }

func main() {
	go Hub.run()
	// http.HandleFunc("/", serverIndex)
	http.HandleFunc("/webSocket", serverWs)
	fmt.Println("Start Server")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Server on")
	}
}
