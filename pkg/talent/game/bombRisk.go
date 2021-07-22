package game

import (
	"encoding/json"
	"math/rand"
	"time"
)

//a="minimum_effort_airport"
//minimum_effort_airport_target
//https://talent.default.serverless-test.moblab-us.cn/api/1/startGame/603f59df15a8cc001645a1bd/backpack

//BombRisk pushpll
type BombRisk struct {
	Delay        int
	currentRound int
}

func NewBombRisk() *BombRisk {
	return &BombRisk{Delay: 8, currentRound: -1}
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

func (hf *BombRisk) GameEnded(g *GameClient, msg *JoginedMsg) bool {
	return false
}

//UserJoined aa
func (hf *BombRisk) UserJoined(g *GameClient, msg *JoginedMsg) {

}

//SessionEnded ss
func (hf *BombRisk) SessionEnded(g *GameClient, msg *SessionEndedMsg) {

}

//GameStated ss
func (hf *BombRisk) GameStated(g *GameClient, mgs *GameStartedMsg) {

}

//GameRoundStarted ss
func (hf *BombRisk) GameRoundStarted(g *GameClient, mgs *GameRoundMsg) {

}

//GameRoundEnded ss
func (hf *BombRisk) GameRoundEnded(g *GameClient, mgs *GameRoundMsg) {
}

//PlayerUpdated ss PLAYER_UPDATED
func (hf *BombRisk) PlayerUpdated(g *GameClient, msg []byte) {
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
			g.stopWatch.End(CHOOSE, "")
			action := Action{Action: CHOOSE, Player: playerUpdated.PlayerNumber}
			rand.Seed(time.Now().UnixNano())
			action.Data = []int{50}
			g.SendActionDelay(action, g.Round, rand.Intn(5)+5)
			return
		}
	}
}
