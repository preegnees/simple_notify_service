package ws

import (
	"log"
	"net/http"

	"notifier/pkg/models"
	pusher "notifier/pkg/push"

	"github.com/gorilla/websocket"
)

type connection struct {}

var storage =  make(map[string]*websocket.Conn)


var _ pusher.IConnection = (*connection)(nil)

var upgrader = websocket.Upgrader{} // use default options

func (c *connection) SendMessage(message models.Message) error { // сюда должно приходить ДТО
	conn, ok := storage[message.Username]
	if !ok {
		return nil
	}
	conn.WriteJSON(message) // исправить потом
	return nil
}

func (c *connection) notif(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	token := r.Header.Get("authorization") // get jwt token, valid and useranme 
	username := token
	storage[username] = conn
}

func New() *connection {

	return &connection{}
}

func (c *connection) Listen() {
	http.HandleFunc("/notif", c.notif)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
