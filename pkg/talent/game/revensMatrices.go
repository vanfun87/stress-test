package game

import (
	"encoding/json"
)

const (
	SOLVE   = "SOLVE"
	RESPOND = "RESPOND"
)

//RevensMatrices Heartflower
type RevensMatrices struct {
	Delay        int
	currentRound int
}

func NewRevensMatrices(delay int) *RevensMatrices {
	return &RevensMatrices{Delay: delay, currentRound: -1}
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
type revenSOLVE struct {
	PlayerUpdated
	Difficulty int           `json:"difficulty"`
	Matrix     [][][]float32 `json:"matrix"`
	Options    [][]float32   `json:"options"`
	Solution   int           `json:"solution"`
}

type revenRESPOND struct {
	PlayerUpdated
	Question  string `json:"question"`
	RangeHigh int    `json:"range_high"`
	RangeLow  int    `json:"range_low"`
}

//UserJoined aa
func (hf *RevensMatrices) UserJoined(g *GameClient, msg *JoinedMsg) {

}

//SessionEnded ss
func (hf *RevensMatrices) SessionEnded(g *GameClient, msg *SessionEndedMsg) {

}

//GameStated ss
func (hf *RevensMatrices) GameStated(g *GameClient, mgs *GameStartedMsg) {

}

//GameRoundStarted ss
func (hf *RevensMatrices) GameRoundStarted(g *GameClient, mgs *GameRoundMsg) {

}

//GameRoundEnded ss
func (hf *RevensMatrices) GameRoundEnded(g *GameClient, mgs *GameRoundMsg) {
	g.stopWatch.End(SOLVE, GAME_ROUND_ENDED)
	g.stopWatch.End(RESPOND, GAME_ROUND_ENDED)
}

//PlayerUpdated ss PLAYER_UPDATED
func (hf *RevensMatrices) PlayerUpdated(g *GameClient, msg []byte) {
	playerUpdated := &PlayerUpdated{}
	err := json.Unmarshal(msg, playerUpdated)
	if err != nil {
		g.stopWatch.Log("json Unmarshal err", err.Error())
		return
	}
	//log.Println(playerUpdated)
	for _, move := range playerUpdated.Moves {
		switch move {
		case SOLVE:
			if hf.currentRound != g.Round {
				hf.currentRound = g.Round
				g.stopWatch.Get(GAME_ROUND_STARTED, SOLVE)
			} else {
				g.stopWatch.End(SOLVE, "")
			}
			var solve revenSOLVE
			err := json.Unmarshal(msg, &solve)
			if err != nil {
				g.stopWatch.Log("json Unmarshal err", err.Error())
				return
			}
			action := Action{Action: SOLVE, Player: playerUpdated.PlayerNumber}
			action.Data = []int{g.Round, solve.Index, solve.Solution}
			g.SendActionDelay(action, g.Round, 5)

			return
		case RESPOND:
			if hf.currentRound != g.Round {
				hf.currentRound = g.Round
				g.stopWatch.Get(GAME_ROUND_STARTED, RESPOND)
			} else {
				g.stopWatch.End(RESPOND, "")
			}
			action := Action{Action: RESPOND, Player: playerUpdated.PlayerNumber}
			action.Data = []int{g.Round, playerUpdated.Index, 3}
			g.SendActionDelay(action, g.Round, 5)

			return
		}

	}
}
