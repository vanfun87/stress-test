package game

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/ginkgoch/stress-test/pkg/talent/lib"

	//"test/websocket"
	"github.com/gorilla/websocket"
)

// The WebsocketClient type represents a game WebSocket connection and game data.
type WebsocketClient struct {
	websocket *websocket.Conn
	messageID int

	serverURL       string
	clientID        string
	sendMsgChan     chan setMessage
	ReceivedMsgChan chan *DataRecv
	handshakeChan   chan error
	closed          bool
	stopWatch       *lib.StopWatch
	userID          int
	sendMode        Mode
	mutex           sync.Mutex
}

type Mode int32

const (
	Single Mode = 1
	Pool   Mode = 2
)

type setMessage interface {
	setMessage(id string, clientID string)
}

//DataSend data need send
type DataSend struct {
	ID       string      `json:"id"`
	Channel  string      `json:"channel"`
	Data     interface{} `json:"data"`
	ClientID string      `json:"clientId"`
}

func (ds *DataSend) setMessage(id string, clientID string) {
	ds.ID = id
	ds.ClientID = clientID
}

type DataRecv struct {
	Channel string
	Data    []byte
}

// Heartbeat represents Heartbeat message
type Heartbeat struct {
	ID             string      `json:"id"`
	Channel        string      `json:"channel"`
	ClientID       string      `json:"clientId"`
	ConnectionType string      `json:"connectionType,omitempty"`
	Ext            interface{} `json:"ext"`
}

func (hb *Heartbeat) setMessage(id string, clientID string) {
	hb.ID = id
	hb.ClientID = clientID
}

//NewGameClient  new a WebsocketClient
func NewWebsocketClient(serverURL string, userid int) *WebsocketClient {
	stopwatch := lib.NewStopWatch(strconv.Itoa(userid))
	gc := WebsocketClient{
		serverURL:       serverURL,
		sendMsgChan:     make(chan setMessage, 10),
		ReceivedMsgChan: make(chan *DataRecv, 10),
		handshakeChan:   make(chan error),
		stopWatch:       &stopwatch,
		userID:          userid,
		sendMode:        Single,
		closed:          true,
	}
	return &gc
}

//SetGameHandler ss
func (ws *WebsocketClient) SetGameHandler() {

}

func (ws *WebsocketClient) sendLoop() {
	ws.stopWatch.Start("sendLoop", "ws.sendMode")
	for sendData := range ws.sendMsgChan {

		if err := ws.sendToWebsocket(sendData); err != nil {
			return
		}
	}
}

func (ws *WebsocketClient) sendData(sendData setMessage) {
	if ws.sendMode == Pool {
		lib.SendWorkPool <- func() {
			if err := ws.sendToWebsocket(sendData); err != nil {
				log.Println("sendData-error:", err)
			}
		}
	} else {
		ws.sendMsgChan <- sendData
	}
}

func (ws *WebsocketClient) sendToWebsocket(sendData setMessage) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	if ws.closed {
		return nil
	}
	ws.messageID++
	sendData.setMessage(strconv.Itoa(ws.messageID), ws.clientID)
	err := ws.websocket.WriteJSON([]setMessage{sendData})
	if err != nil {
		return err
	}
	return nil
}

//Run game
func (ws *WebsocketClient) Connect() error {
	err := ws.connect()
	if err != nil {
		return err
	}
	if ws.sendMode == Single {
		go ws.sendLoop()
	}

	go ws.handleMessage()
	if err = <-ws.handshakeChan; err != nil {
		ws.close()
	}

	return err
}

//Connect to game server
func (ws *WebsocketClient) connect() error {
	ws.stopWatch.Start("connect", ws.serverURL)
	c, _, err := websocket.DefaultDialer.Dial(ws.serverURL, nil)
	ws.stopWatch.End("connect", fmt.Sprintf("%v", err))
	if err != nil {
		log.Println("Dial error ", err)
		return err
	}
	ws.websocket = c
	ws.closed = false
	return ws.handshake()
}

func (ws *WebsocketClient) handshake() error {
	var handshakeBody = `[
	{
		"id":"1",
		"version":"1.0",
		"minimumVersion":"1.0",
		"channel":"/meta/handshake",
		"supportedConnectionTypes":[
			"websocket",
			"long-polling",
			"callback-polling"
		],
		"advice":{
			"timeout":60000,
			"interval":0
		},
		"ext":{
			"ack":true
		}
	  }
  ]`

	ws.stopWatch.Start("handshake", "")
	ws.stopWatch.Log("send handshake", handshakeBody)
	err := ws.websocket.WriteMessage(websocket.TextMessage, []byte(handshakeBody))

	if err != nil {
		log.Println("write error:", err)
		return err
	}
	return nil
}

func (ws *WebsocketClient) handleMessage() {
	for {
		_, message, err := ws.websocket.ReadMessage()
		if err != nil {
			if !ws.closed {
				ws.stopWatch.Log("error", err.Error())
				ws.ReceivedMsgChan <- &DataRecv{Channel: "error", Data: []byte(err.Error())}
			}
			close(ws.ReceivedMsgChan)
			return
		}
		//ws.stopWatch.Log("recv msg:", string(message))
		var rsMsgs []ReceivedMsg
		//3.解析
		err = json.Unmarshal(message, &rsMsgs)
		if err != nil {
			ws.stopWatch.Log("json error", err.Error())
			continue
		}
		if rsMsgs[0].Error != "" {
			ws.stopWatch.Log("rsMsg error", rsMsgs[0].Error)
			ws.close()
		}
		switch rsMsg := rsMsgs[0]; rsMsg.Channel {
		case "/meta/handshake":
			ws.stopWatch.End("handshake", "")
			err = ws.handleHandshake(message)
			ws.handshakeChan <- err

		case "/meta/connect":
			ws.handleHeartbeat(message)
		default:
			fmt.Println("rsMsg.Data", rsMsg.Data)
			data, err := json.Marshal(rsMsg.Data)
			if err != nil {
				ws.stopWatch.Log("json error", err.Error())
			} else {
				ws.ReceivedMsgChan <- &DataRecv{Channel: rsMsg.Channel, Data: data}
			}
		}

	}

}

func (ws *WebsocketClient) close() {
	ws.closed = true
	if ws.websocket != nil {
		ws.websocket.Close()
	}

}

//SendAction send action, use go SendAction
func (ws *WebsocketClient) SendAction(action interface{}, channel string) {
	msgNeedSend := DataSend{
		Data:    action,
		Channel: channel,
	}
	ws.sendData(&msgNeedSend)

}

func (ws *WebsocketClient) handleHeartbeat(message []byte) {
	var heartbeats []Heartbeat
	err := json.Unmarshal(message, &heartbeats)
	if err != nil {
		log.Println("handshakeMsg json err:", err, heartbeats)
		return
	}
	heartbeatMsg := heartbeats[0]
	heartbeatMsg.ConnectionType = "websocket" //connectionType: 'websocket'
	ws.sendData(&heartbeatMsg)
}

func (ws *WebsocketClient) handleHandshake(msg []byte) error {
	log.Println("handshake recv")
	var handshakeMessages []HandshakeMsg
	err := json.Unmarshal(msg, &handshakeMessages)
	if err != nil {
		log.Println("handshakeMsg json err:", err, handshakeMessages)
		return err
	}
	handshakeMsg := handshakeMessages[0]
	ws.clientID = handshakeMsg.ClientID
	heartbeatBody := &Heartbeat{
		ConnectionType: "websocket",
		Ext:            map[string]interface{}{"ack": 0},
		Channel:        "/meta/connect",
	}
	ws.sendData(heartbeatBody)
	return nil
}
