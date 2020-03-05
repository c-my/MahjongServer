package game

import (
	"github.com/c-my/MahjongServer/message"
	"github.com/c-my/MahjongServer/model"
	"github.com/c-my/MahjongServer/rule"
	"log"
)

type MahjongManager struct {
	wall       *model.Wall
	playerTile []model.PlayerTile
	rules      rule.MahjongRule

	gameRecvCh chan message.GameMsgRecv
	gameSendCh chan message.GameMsgSend
}

//NewMahjongManager return new Mahjong instance, with initiated game channels
// and tile-wall
func NewMahjongManager(gameRecvCh chan message.GameMsgRecv, gameSendCh chan message.GameMsgSend, r rule.MahjongRule) *MahjongManager {
	mahjong := MahjongManager{}
	mahjong.rules = r
	mahjong.wall = model.NewWall(r.GenerateTiles)
	mahjong.playerTile = make([]model.PlayerTile, 4)
	mahjong.gameRecvCh = gameRecvCh
	mahjong.gameSendCh = gameSendCh
	return &mahjong
}

//Start will start the main loop of a game in a new goroutine
func (m *MahjongManager) Start() {
	go m.gameLoop(m.gameRecvCh, m.gameSendCh)
}

func (m MahjongManager) GetFirstPlayer() {
	panic("implement me")
}

func (m *MahjongManager) gameLoop(gameRecvCh chan message.GameMsgRecv, gameSendCh chan message.GameMsgSend) {
	for {
		select {
		case msg := <-gameRecvCh:
			log.Println("receive from game loop: ", msg)
			switch msg.Action {
			case message.Start:
				m.dealTile()
				availableActions := make([]int, 0)
				availableActions = append(availableActions, message.Deal)
				tilesCount := make([]int, 0)
				for i := 0; i < 4; i++ {
					tilesCount = append(tilesCount, len(m.playerTile[i].HandTiles))
				}
				var msgSend message.GameMsgSend
				msgSend.MsgType = message.GameMsgType
				msgSend.AvailableActions = availableActions
				msgSend.CurrentTile = m.playerTile[0].HandTiles
				msgSend.TilesCount = tilesCount
				gameSendCh <- msgSend
			case message.Discard:

			}
			// TODO:main logic here

		}
	}
}

func (m *MahjongManager) onStart(msg message.GameMsgRecv) {
	m.dealTile()

	var msgSend message.GameMsgSend
	msgSend.MsgType = message.GameMsgType

}

//func (m *MahjongManager) getAvailableActions(tiles []model.Tile, newTile model.Tile) []int {
//	availableActions := make([]int, 0)
//	if m.rules.CanChow(tiles, newTile) {
//		availableActions = append(availableActions, message.Chow)
//	}
//	if m.rules.CanPong(tiles, newTile) {
//		availableActions = append(availableActions, message.Pong)
//	}
//	if m.rules.CanKong(tiles, newTile) {
//		availableActions = append(availableActions, message.Kong)
//	}
//	if m.rules.CanWin(tiles, newTile) {
//		availableActions = append(availableActions, message.Win)
//	}
//	return availableActions
//}

func (m *MahjongManager) dealTile() {
	for i := 0; i < 3; i++ {
		m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, m.wall.FrontDraw())
		m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, m.wall.FrontDraw())
		m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, m.wall.FrontDraw())
		m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, m.wall.FrontDraw())
		m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, m.wall.FrontDraw())
		m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, m.wall.FrontDraw())
		m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, m.wall.FrontDraw())
		m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, m.wall.FrontDraw())
		m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, m.wall.FrontDraw())
		m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, m.wall.FrontDraw())
		m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, m.wall.FrontDraw())
		m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, m.wall.FrontDraw())
		m.playerTile[3].HandTiles = append(m.playerTile[3].HandTiles, m.wall.FrontDraw())
		m.playerTile[3].HandTiles = append(m.playerTile[3].HandTiles, m.wall.FrontDraw())
		m.playerTile[3].HandTiles = append(m.playerTile[3].HandTiles, m.wall.FrontDraw())
		m.playerTile[3].HandTiles = append(m.playerTile[3].HandTiles, m.wall.FrontDraw())
	}
	m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, m.wall.FrontDraw())
	m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, m.wall.FrontDraw())
	m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, m.wall.FrontDraw())
	m.playerTile[3].HandTiles = append(m.playerTile[3].HandTiles, m.wall.FrontDraw())

	for i := 0; i <= 3; i++ {
		m.playerTile[i].SortHand()
	}
}
