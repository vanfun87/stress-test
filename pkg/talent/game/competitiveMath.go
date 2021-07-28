package game

import (
	"encoding/json"
)

const (
	CONFIDENCE      = "CONFIDENCE"
	RANK            = "RANK"
	ANSWER          = "ANSWER"
	RANK_CONFIDENCE = "RANK_CONFIDENCE"
	CHOOSE          = "CHOOSE"
	FINISH_TUTORIAL = "FINISH_TUTORIAL"
)

//CompetitiveMath Competitive Math
type CompetitiveMath struct {
	Delay        int
	currentRound int
	Player       CompetitiveMathPlayer
}

func NewCompetitiveMath(delay int) *CompetitiveMath {
	defaultPlayer := CompetitiveMathPlayer{
		CONFIDENCE:     10, // answer number
		BetRatio:       1,  //ratio
		RANK:           1,
		RankConfidence: 7,
		CHOOSE:         "single",
		Answer: [8][]int{
			{}, //confidence
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{}, //finish TUTORIAL
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{}, //RANK,RANK CONFIDENCE
			{}, //CHOOSE
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	return &CompetitiveMath{Delay: delay, currentRound: -1, Player: defaultPlayer}
}

//UserJoined aa
func (hf *CompetitiveMath) UserJoined(g *GameClient, msg *JoinedMsg) {

}

//SessionEnded ss
func (hf *CompetitiveMath) SessionEnded(g *GameClient, msg *SessionEndedMsg) {
}

//GameStated ss
func (hf *CompetitiveMath) GameStated(g *GameClient, mgs *GameStartedMsg) {

}

//GameRoundStarted ss
func (hf *CompetitiveMath) GameRoundStarted(g *GameClient, mgs *GameRoundMsg) {

}

//GameRoundEnded ss
func (hf *CompetitiveMath) GameRoundEnded(g *GameClient, mgs *GameRoundMsg) {
	g.stopWatch.End(ANSWER, GAME_ROUND_ENDED)
}

type CompetitiveMathPlayer struct {
	CONFIDENCE     int
	Answer         [8][]int
	BetRatio       int
	RANK           int
	RankConfidence int
	CHOOSE         string
	answercount    int
}

type CompetitiveMathAnswer struct {
	PlayerNumber        int           `json:"playerNumber"`
	RoundCorrectAnswers []interface{} `json:"round_correct_answers"`
	CurrIndex           int           `json:"curr_index"`
	Moves               []string      `json:"moves"`
	RoundAnswers        []interface{} `json:"round_answers"`
	Numbers             []int         `json:"numbers"`
	RoundPayoff         int           `json:"round_payoff"`
	Choice              string        `json:"choice"`
	Earning             int           `json:"earning"`
	Target              int           `json:"target"`
}

//PlayerUpdated ss PLAYER_UPDATED
func (hf *CompetitiveMath) PlayerUpdated(g *GameClient, msg []byte) {
	playerUpdated := &PlayerUpdated{}
	err := json.Unmarshal(msg, playerUpdated)
	if err != nil {
		g.stopWatch.Log("json Unmarshal err", err.Error())
		return
	}

	for _, move := range playerUpdated.Moves {
		switch move {
		case ANSWER:
			if hf.currentRound != g.Round {
				hf.currentRound = g.Round
				hf.Player.answercount = 0
				g.stopWatch.Get(GAME_ROUND_STARTED, ANSWER)
			} else {
				g.stopWatch.End(ANSWER, "")
			}
			var answer CompetitiveMathAnswer
			err := json.Unmarshal(msg, &answer)
			if err != nil {
				g.stopWatch.Log("json Unmarshal err", err.Error())
				return
			}
			if g.Round > len(hf.Player.Answer) {
				//g.stopWatch.Log("receive  answer action in g.Round ", fmt.Sprintf("%d", g.Round))
				return
			}
			if len(hf.Player.Answer[g.Round]) <= hf.Player.answercount {
				return
			}

			action := Action{Action: ANSWER, Player: playerUpdated.PlayerNumber}
			playeranswer := getCompetitiveMathRightAnswer(answer.Numbers, answer.Target)
			action.Data = append([]int{g.Round, answer.CurrIndex}, playeranswer[:]...) //[round,confidence]

			hf.Player.answercount++
			g.SendActionDelay(action, g.Round, 5)
			return
		case RANK:
			g.stopWatch.Get(GAME_ROUND_STARTED, RANK)
			action := Action{Action: RANK, Player: playerUpdated.PlayerNumber}
			action.Data = []int{g.Round, hf.Player.RANK} //[round,confidence]
			g.SendActionDelay(action, g.Round, 5)
			return
		case CONFIDENCE:
			g.stopWatch.Get(GAME_ROUND_STARTED, CONFIDENCE)
			action := Action{Action: CONFIDENCE, Player: playerUpdated.PlayerNumber}
			action.Data = []int{g.Round, hf.Player.CONFIDENCE} //[round,confidence]
			g.SendActionDelay(action, g.Round, 5)
			return

		case RANK_CONFIDENCE:
			g.stopWatch.Get(GAME_ROUND_STARTED, RANK_CONFIDENCE)
			action := Action{Action: RANK_CONFIDENCE, Player: playerUpdated.PlayerNumber}
			action.Data = []int{g.Round, hf.Player.CONFIDENCE} //[round,confidence]
			g.SendActionDelay(action, g.Round, 5)
			return

		case CHOOSE:
			g.stopWatch.Get(GAME_ROUND_STARTED, CHOOSE)
			action := Action{Action: CHOOSE, Player: playerUpdated.PlayerNumber}
			var data []interface{}
			data = append(data, g.Round)
			data = append(data, hf.Player.CHOOSE)
			action.Data = data //[round,choose]
			g.SendActionDelay(action, g.Round, 5)
			return
		case FINISH_TUTORIAL:
			g.stopWatch.End(ANSWER, FINISH_TUTORIAL)
			action := Action{Action: FINISH_TUTORIAL, Player: playerUpdated.PlayerNumber}
			action.Data = []int{g.Round}
			g.SendActionDelay(action, g.Round, 5)
			return

		}

	}
}

func getCompetitiveMathRightAnswer(data []int, target int) (answer [2]int) {
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i]+data[j] == target {
				answer[0] = i
				answer[1] = j
				return
			}
		}
	}
	answer[1] = 1
	return
}
