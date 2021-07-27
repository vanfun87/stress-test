package game

import (
	"encoding/json"
	"math/rand"
	"time"
)

//PushPull pushpll
type PushPull struct {
	Delay        int
	currentRound int
}

func NewPushPull() *PushPull {
	return &PushPull{Delay: 1, currentRound: -1}
}

// {
// 	"playerNumber": 0,
// 	"difficulty": 0,
// 	"period": 1,
// 	"solution": 1,
// 	"moves": ["SOLVE"],
// 	"options": [
// 		[4, 0.8, 30, 2],
// 		[0, 1.0, 30, 0],
// 		[10, 0.8, 180, 0],
// 		[7, 1.0, 150, 0],
// 		[10, 1.0, 30, 3],
// 		[0, 0.8, 90, 2]
// 	],
// 	"index": 0,
// 	"matrix": [
// 		[
// 			[0, 1.0, 30, 0],
// 			[0, 1.0, 30, 0],
// 			[0, 1.0, 30, 0]
// 		],
// 		[
// 			[0, 1.0, 30, 0],
// 			[0, 1.0, 30, 0],
// 			[0, 1.0, 30, 0]
// 		],
// 		[
// 			[0, 1.0, 30, 0],
// 			[0, 1.0, 30, 0], null
// 		]
// 	],
// 	"earning": 0
// }

//UserJoined aa
func (hf *PushPull) UserJoined(g *GameClient, msg *JoinedMsg) {

}

//SessionEnded ss
func (hf *PushPull) SessionEnded(g *GameClient, msg *SessionEndedMsg) {

}

//GameStated ss
func (hf *PushPull) GameStated(g *GameClient, mgs *GameStartedMsg) {

}

//GameRoundStarted ss
func (hf *PushPull) GameRoundStarted(g *GameClient, mgs *GameRoundMsg) {

}

//GameRoundEnded ss
func (hf *PushPull) GameRoundEnded(g *GameClient, mgs *GameRoundMsg) {
	g.stopWatch.End(CHOOSE, GAME_ROUND_ENDED)
}

//PlayerUpdated ss PLAYER_UPDATED
func (hf *PushPull) PlayerUpdated(g *GameClient, msg []byte) {
	playerUpdated := &PlayerUpdated{}
	err := json.Unmarshal(msg, playerUpdated)
	if err != nil {
		g.stopWatch.Log("json Unmarshal err", err.Error())
		return
	}
	//log.Println(playerUpdated)
	for _, move := range playerUpdated.Moves {
		switch move {
		case CHOOSE:
			if hf.currentRound != g.Round {
				hf.currentRound = g.Round
				g.stopWatch.Get(GAME_ROUND_STARTED, CHOOSE)
			} else {
				return
			}
			action := Action{Action: CHOOSE, Player: playerUpdated.PlayerNumber}
			rand.Seed(time.Now().UnixNano())
			p := rand.Intn(2)
			action.Data = []int{g.Round, p}
			g.SendActionDelay(action, g.Round, 1)

			return
		}
	}
}
