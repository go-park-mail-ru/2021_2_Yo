package websocket

import (
	log "backend/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ID struct {
	//TODO: Скорее всего, здесь нужен тег
	ID string `json:"id,omitempty"`
}

func GetID(conn *websocket.Conn) (string, error) {
	/*var userID ID
	log.Info("here in id")
	err := conn.ReadJSON(userID)
	log.Error(err)
	if err != nil {
		return "", err
	}
	return userID.ID, nil
	*/
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Error("read message", err)
	}
	uId := ID{}
	err = json.Unmarshal(p, &uId)
	if err != nil {
		log.Error("unmarshal",err)
	}
	return uId.ID, nil
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("err = ", err)
		return nil, err
	}
	return conn, nil
}
