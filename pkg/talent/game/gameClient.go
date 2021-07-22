package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ginkgoch/stress-test/pkg/talent/lib"
)

// The GameClient type represents a game WebSocket connection and game data.
type GameClient struct {
	userID     int
	roomID     string
	gamePlayer GamePlayer
	Round      int
	gameID     string
	wsClient   *WebsocketClient
	stopWatch  *lib.StopWatch
	Delay      int
}

////NewGameClient  new a GameClient
//func NewGameClient(userID int, roomID string, serverURL string, player GamePlayer, phoneNumber string) *GameClient {
//	stopwatch := lib.NewStopWatch(phoneNumber + ":" + strconv.Itoa(userID))
//	gc := GameClient{
//		userID:     userID,
//		roomID:     roomID,
//		gamePlayer: player,
//		wsClient:   NewWebsocketClient(serverURL, userID),
//		stopWatch:  &stopwatch,
//	}
//	return &gc
//}

//NewGameClient  new a GameClient
func NewGameClient(config *GameConfig, player GamePlayer) *GameClient {
	stopwatch := lib.NewStopWatch(config.PhoneNumber + ":" + strconv.Itoa(config.PlayerID))

	wsUrl := fmt.Sprintf("wss://%s/game-server/cometd", config.WebSocketHost())
	fmt.Println("ws server", wsUrl)
	gc := GameClient{
		userID:     config.PlayerID,
		roomID:     config.RoomID,
		gamePlayer: player,
		wsClient:   NewWebsocketClient(wsUrl, config.PlayerID),
		stopWatch:  &stopwatch,
	}
	return &gc
}

//Run game
func (g *GameClient) Run() (err error) {

	defer g.close()
	g.stopWatch.Start("Run", strconv.Itoa(g.userID))
	defer func() {
		g.stopWatch.End("Run", fmt.Sprintf("%d %v", g.userID, err))
	}()
	g.stopWatch.Start("connect_wsclient", "")
	err = g.wsClient.Connect()
	g.stopWatch.End("connect_wsclient", fmt.Sprintf("%v", err))
	if err != nil {
		return
	}
	g.joinGame()
	g.stopWatch.Start("startJoin", "")

	fmt.Println("starting handle message")

	time.Sleep(10 * time.Second)
	err = g.handleMessage()
	g.stopWatch.Log("DelayTime:", strconv.Itoa(g.Delay))
	if err != nil {
		return err
	}
	return nil

}

func (g *GameClient) handleMessage() error {
	// select{
	// case _, message, err := g.websocket.ReadMessage():
	// 	fmt.Println("sefe")
	// }

	for receivMsg := range g.wsClient.ReceivedMsgChan {
		fmt.Println(receivMsg)

		switch receivMsg.Channel {
		case "error":
			return errors.New(string(receivMsg.Data))
		case "/gameroom":
			event := Event{}
			err := json.Unmarshal(receivMsg.Data, &event)
			if err != nil {
				return err
			}
			switch event.Event {
			case UNAVAILABLE:
				return errors.New("game UNAVAILABLE")
			case USER_JOINED:
				g.stopWatch.End("startJoin", USER_JOINED)
				joginedMsg := &JoginedMsg{}
				err = json.Unmarshal(receivMsg.Data, joginedMsg)
				if err != nil {
					return err
				}
				if !joginedMsg.Active {
					return errors.New("join game ,not active")
				}
				g.stopWatch.Start(USER_JOINED, "")
				g.gamePlayer.UserJoined(g, joginedMsg)
			case SESSION_ENDED:
				g.stopWatch.End(USER_JOINED, SESSION_ENDED)
				g.stopWatch.End(GAME_ENDED, SESSION_ENDED)
				joginedMsg := &SessionEndedMsg{}
				err = json.Unmarshal(receivMsg.Data, joginedMsg)
				if err != nil {
					return err
				}
				g.gamePlayer.SessionEnded(g, joginedMsg)
				return nil
			default:
				g.stopWatch.Start("/gameroom unhandle event: ", event.Event)
			}
		case "/game":
			event := Event{}
			err := json.Unmarshal(receivMsg.Data, &event)
			eventData, err := json.Marshal(event.Data)
			if err != nil {
				g.stopWatch.Log("json Unmarshal error", err.Error())
			}
			switch event.Event {
			case GAME_STARTED:
				g.stopWatch.Start(GAME_STARTED, "")
				g.stopWatch.Start(GAME_ROUND_STARTED, "FAKE")
				msg := &GameStartedMsg{}
				err = json.Unmarshal(eventData, msg)
				if err != nil {
					g.stopWatch.Log("json Unmarshal error", err.Error())
					return err
				}
				if msg.Status != "RUNNING" {
					g.stopWatch.Log("game status", msg.Status)
					return nil
				}
				g.gameID = msg.GameID
				g.Round = msg.Round
				g.gamePlayer.GameStated(g, msg)
			case PLAYER_UPDATED:
				g.playerUpdated(eventData)

			case GAME_ROUND_STARTED:
				g.stopWatch.Start(GAME_ROUND_STARTED, "")
				g.stopWatch.End(GAME_ROUND_ENDED, GAME_ROUND_STARTED+":"+strconv.Itoa(g.Round))
				msg := &GameRoundMsg{}
				err = json.Unmarshal(eventData, msg)
				if err != nil {
					g.stopWatch.Log("json Unmarshal error", err.Error())
					return err
				}
				g.Round = msg.Round
				g.gamePlayer.GameRoundStarted(g, msg)
			case GAME_ROUND_ENDED:
				g.stopWatch.End(GAME_ROUND_STARTED, GAME_ROUND_ENDED+":"+strconv.Itoa(g.Round))
				g.stopWatch.Start(GAME_ROUND_ENDED, "")
				msg := &GameRoundMsg{}
				err = json.Unmarshal(eventData, msg)
				if err != nil {
					g.stopWatch.Log("json Unmarshal error", err.Error())
					return err
				}
				g.Round = msg.Round
				g.gamePlayer.GameRoundEnded(g, msg)
			case GAME_ENDED:
				g.stopWatch.End(GAME_STARTED, GAME_ENDED)
				g.stopWatch.Start(GAME_ENDED, "")
				return nil
			default:
				g.stopWatch.Log("/game unhandle event: ", event.Event)
			}
		default:
			g.stopWatch.Log("unhandle_chnnel", receivMsg.Channel)
		}
	}
	return nil
}

func (g *GameClient) playerUpdated(eventData []byte) {
	playerUpdated := &PlayerUpdated{}
	err := json.Unmarshal(eventData, playerUpdated)
	if err != nil {
		g.stopWatch.Log("error", err.Error())
		return
	}
	g.stopWatch.Log("PLAYER_UPDATED_moves", strings.Join(playerUpdated.Moves, ","))

	g.stopWatch.Start(PLAYER_UPDATED, "round:"+strconv.Itoa(g.Round))
	g.gamePlayer.PlayerUpdated(g, eventData)
}

func (g *GameClient) close() {
	g.wsClient.close()
}

//SendActionDelay send action,  delay second
func (g *GameClient) SendActionDelay(action Action, round int, delay int) {
	g.Delay += delay
	lib.SendWork(func() {
		g.SendAction(action, round)
	}, delay)
}

//SendAction send action, use go SendAction
func (g *GameClient) SendAction(action Action, round int) {
	if g.Round != round { // need lock?  g.Round
		return
	}
	channel := "/service/game/" + g.gameID
	g.stopWatch.Start(action.Action, "")
	g.wsClient.SendAction(action, channel)
}

func (g *GameClient) joinGame() {
	joinGame := JoinGameSend{
		Action: "join",
		Room:   g.roomID,
		User:   g.userID,
	}
	g.wsClient.SendAction(joinGame, "/service/gameroom/"+g.roomID)
}

type GameID string

const (
	RM    GameID = "ravens_matrices"
	CM    GameID = "competitive_math"
	PP    GameID = "push_pull"
	AIR   GameID = "minimum_effort_airport"
	AIR_T GameID = "minimum_effort_airport_target"
	//a="minimum_effort_airport"
	//minimum_effort_airport_target
	//HF              GameID = "hearts_flowers"
	//CB              GameID = "Corsiblocks"
	//AC              GameID = "attention_control"
	//NUMERACY        GameID = "numeracy"
	//GRIT_ASSESSMENT GameID = "grit_assessment"
)

func RunGame(gamconfig *GameConfig) (err error) {
	var gp GamePlayer
	switch GameID(gamconfig.ID) {
	case RM:
		gp = NewRevensMatrices(5)
	case CM:
		gp = NewCompetitiveMath(5)
	case PP:
		gp = NewPushPull()
	case AIR:
		gp = NewAirport()
	case AIR_T:
		gp = NewAirport()
	default:
		err = fmt.Errorf("player %d, no such game:%s ", gamconfig.PlayerID, gamconfig.ID)
		return
	}
	gc := *NewGameClient(gamconfig, gp)
	err = gc.Run()
	return
}
