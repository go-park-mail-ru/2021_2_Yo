package websocket

import (
	log "backend/pkg/logger"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ID struct {
	//TODO: Скорее всего, здесь нужен тег
	ID string
}

func GetID(conn *websocket.Conn) (string, error) {
	var userID ID
	err := conn.ReadJSON(userID)
	if err != nil {
		return "", err
	}
	return userID.ID, nil
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("err = ", err)
		return nil, err
	}
	return conn, nil
}


