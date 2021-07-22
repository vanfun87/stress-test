package game

//GamePlayer  gameplayer interface
type GamePlayer interface {
	UserJoined(*GameClient, *JoginedMsg)
	SessionEnded(*GameClient, *SessionEndedMsg)
	GameStated(*GameClient, *GameStartedMsg)
	PlayerUpdated(*GameClient, []byte)
	GameRoundStarted(*GameClient, *GameRoundMsg)
	GameRoundEnded(*GameClient, *GameRoundMsg)
}

// channel:

// internal:
// /meta/handshake  handshake
// /meta/connect   heartbeat

// default_game:
// /gameroom  USER_JOINED,UNAVAILABLE,SESSION_ENDED
// /game       GAME_STARTED,GAME_ENDED, GAME_ROUND_STARTED,GAME_ROUND_ENDED ,PLAYER_UPDATED

//  PLAYER_UPDATED moves:
