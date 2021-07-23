package talent

import (
	"fmt"
	"net/http"

	"github.com/ginkgoch/stress-test/pkg/talent/game"
)

type TalentObject struct {
	ServiceEndpoint string
	UserId          string
	Cookie          *http.Cookie
	GameConfig      *game.GameConfig
}

func (t *TalentObject) String() string {
	s := fmt.Sprintln()
	s += fmt.Sprintln("ServiceEndpoint:", t.ServiceEndpoint)
	s += fmt.Sprintln("UserId:", t.UserId)
	s += fmt.Sprintln("Cookie:", t.Cookie)
	s += fmt.Sprintln("GameConfig:", t.GameConfig)
	return s
}

type Information struct {
	Success bool `json:"success"`
	User    struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		ID          string `json:"id"`
	} `json:"user"`
}

// type GameConfig struct {
// 	ID       string `json:"id"`
// 	PlayerID int    `json:"playerId"`
// 	RoomID   string `json:"roomId"`
// 	Server   string `json:"server"`
// 	GameURL  string `json:"gameurl"`
// }

type StartGameData struct {
	Success bool `json:"success"`
	Data    struct {
		ID       string `json:"id"`
		PlayerID string `json:"playerId"`
		RoomID   string `json:"roomId"`
		Server   string `json:"serverAddress"`
		Gameurl  string `json:"gameurl"`
	} `json:"data"`
}
