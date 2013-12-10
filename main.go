package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    "flag"
    zmq "github.com/pebbe/zmq3"
)

var addr = flag.String("addr", ":8080", "http service address")
var responder *(zmq.Socket)

func main() {
    flag.Parse()

    responder, _ = zmq.NewSocket(zmq.SUB)
    defer socket.Close()
    responder.Connect("tcp://localhost:5563")
    responder.SetSubscribe("message")


    http.HandleFunc("/ws", handler)
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(w, "Not a websocket handshake", 400)
        return
    } else if err != nil {

        log.Println(err)
        return
    }

    log.Println("Start...")
    for {
        msg, _ := responder.RecvMessage(0)
        log.Println("Received ", msg)    
        //websockets
        if err := conn.WriteMessage(websocket.TextMessage, []byte(msg[1])); err != nil {
            log.Println(err)
            return
        }

    }
}