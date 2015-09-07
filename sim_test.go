package simulator

import (
	"bytes"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"

	"golang.org/x/net/websocket"
)

func TestWebSocketConnect(t *testing.T) {
	ts := httptest.NewServer(WebSockerHandler)
	defer ts.Close()

	url := "ws" + ts.URL[4:]

	ws, _ := websocket.Dial(url, "application/json", "http://*")
	defer ws.Close()

	websocket.JSON.Send(ws, map[string]string{
		"event": "connect",
	})
	// DefaultConn.Close()
	DefaultConn = nil
}

func TestRestRequest(t *testing.T) {
	ts := httptest.NewServer(WebSockerHandler)
	defer ts.Close()

	url := "ws" + ts.URL[4:]

	ws, _ := websocket.Dial(url, "application/json", "http://*")
	defer ws.Close()

	websocket.JSON.Send(ws, map[string]string{
		"event": "connect",
	})

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/softphone", SimulatorHandler),
	)
	if err != nil {
		log.Fatal(err)
	}

	input := bytes.NewReader([]byte(`{
        "ani":"66801234567",
        "nccacallheaderid":"1-5432123"
        }`))

	api.SetApp(router)
	recorded := test.RunRequest(t, api.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/softphone", input))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	var data softphone
	websocket.JSON.Receive(ws, &data)
	fmt.Println(data)

	DefaultConn = nil
}
