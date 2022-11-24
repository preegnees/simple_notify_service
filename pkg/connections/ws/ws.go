package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"

	dto "notifier/pkg/dto"
	models "notifier/pkg/models"
	pusher "notifier/pkg/push"

	"github.com/gorilla/websocket"
)

type connection struct {
	addr    string
	storage map[string]*websocket.Conn
	mx      sync.Mutex
}

var _ pusher.IConnection = (*connection)(nil)

func (c *connection) SendMessage(dtoMessage dto.DTOMessagePush) error {
	
	conn, ok := c.storage[dtoMessage.Username]
	if !ok {
		return errors.New("conn did not found")
	}

	payload := models.MessagePush{
		MessageSubject: dtoMessage.MessageSubject,
		Message:        dtoMessage.Message,
	}

	return conn.WriteJSON(payload) // TODO(замаршалить потом)
}

var upgrader = websocket.Upgrader{}

func (c *connection) notif(w http.ResponseWriter, r *http.Request) {

	jwtToken := r.Header.Get("auth_token")
	respHeader := make(http.Header)
	respHeader.Add("signin", "true")
	// if jwtToken != "jwtTokenTest" { // TODO(проверить на соответсвие токен jwt)
	// 	respHeader.Set("signin", "false")
	// }

	conn, err := upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	username := jwtToken // TODO(получить username из токена)
	
	c.mx.Lock()
	defer c.mx.Unlock()
	c.storage[username] = conn
	log.Println(c.storage)
	w.WriteHeader(http.StatusOK)
}

func New(addr string) *connection {

	if addr == "" {
		addr = "localhost:8080"
	}
	return &connection{
		addr:    addr,
		storage: make(map[string]*websocket.Conn),
	}
}

func (c *connection) Listen() error {

	http.HandleFunc("/notif", c.notif)
	return http.ListenAndServe(c.addr, nil)
}

func (c *connection) CloseConn(username string) error {

	conn, ok := c.storage[username]
	if !ok {
		return nil
	}
	return conn.Close()
}
