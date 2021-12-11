package easyWebsocket

import (
    "net/http"
	log "github.com/sirupsen/logrus"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

type ID struct {
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
        log.Println(err)
        return nil, err
    }

    return conn, nil
}

