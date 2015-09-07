package simulator

import (
	"fmt"
	"reflect"

	"github.com/ant0ine/go-json-rest/rest"

	"golang.org/x/net/websocket"
)

// Default Connection
var DefaultConn *websocket.Conn

// WebSockerHandler use as http.Handler
var WebSockerHandler = websocket.Handler(webSockerHandler)

func webSockerHandler(conn *websocket.Conn) {
	DefaultConn = conn
	data := make(chan map[string]string)
	err := make(chan error)
	go receiever(data, err, conn)

	for {
		if DefaultConn == nil {
			break
		}
		select {
		case e := <-err:
			fmt.Println("error:", e)
		case d := <-data:
			fmt.Println("data:", d)
		}
	}

}

func receiever(dataChan chan map[string]string, errChan chan error, conn *websocket.Conn) {
	for {
		if DefaultConn == nil {
			break
		}
		data := map[string]string{}
		err := websocket.JSON.Receive(conn, &data)
		if err != nil && err.Error() != "EOF" {
			errChan <- err
		}
		if !reflect.DeepEqual(data, map[string]string{}) {
			dataChan <- data
		}
	}
}

type softphone struct {
	ANI              string `json:"ani"`
	NCCACallheaderId string `json:"nccacallheaderid"`
}

func SimulatorHandler(w rest.ResponseWriter, r *rest.Request) {
	var sp softphone
	r.DecodeJsonPayload(&sp)

	w.WriteJson(map[string]string{
		"status": "success",
	})
	websocket.JSON.Send(DefaultConn, sp)
}
