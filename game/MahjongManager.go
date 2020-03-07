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

	lastTableOrder    int
	currentTableOrder int

	lastTile    model.Tile
	currentTile model.Tile

	lastAction int
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
		msg := <-gameRecvCh
		log.Printf("player:[%d], action:[%d] ", msg.TableOrder, msg.Action)
		switch msg.Action {
		case message.Start:
			m.handleStart(msg)
		case message.Discard:
			m.handleDiscard(msg)
		case message.Chow:
			m.handleChow(msg)
		case message.Pong:
			m.handlePong(msg)
		case message.Win:
			m.handleWin(msg)
		case message.Cancel:

		}
		// TODO:main logic here
	}
}

func (m *MahjongManager) handleStart(msg message.GameMsgRecv) {
	m.dealTile()
	availableActions := make([]int, 0)
	availableActions = append(availableActions, message.Deal)

	var msgSend message.GameMsgSend
	msgSend.MsgType = message.GameMsgType
	msgSend.AvailableActions = availableActions
	msgSend.CurrentTile = m.playerTile[0].HandTiles
	msgSend.PlayerTile = m.playerTile
	msgSend.WallCount = m.wall.Length()
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) Draw(tableOrder int) {
	player := m.playerTile[tableOrder]
	newTile := m.wall.FrontDraw()
	currentTile := make([]model.Tile, 0)
	currentTile = append(currentTile, newTile)
	var msg message.GameMsgSend
	msg.MsgType = message.GameMsgType
	msg.AvailableActions = m.getAvailableActions(player.HandTiles, player.ShownTiles, newTile)
	msg.CurrentTile = currentTile
}

func (m *MahjongManager) handleDiscard(msg message.GameMsgRecv) {
	currentOrder := msg.TableOrder
	// 检查是否有人胡牌
	for i := 1; i <= 3; i++ {
		hand := m.playerTile[(currentOrder+i)%4].HandTiles
		shown := m.playerTile[(currentOrder+i)%4].ShownTiles
		if m.rules.CanWin(hand, shown, msg.Tile) {
			//TODO:send msg to this potential winner
			return
		}
	}
	// 检查是否有人碰、杠
	for i := 1; i <= 3; i++ {
		hand := m.playerTile[(currentOrder+i)%4].HandTiles
		shown := m.playerTile[(currentOrder+i)%4].ShownTiles
		if m.rules.CanPong(hand, msg.Tile) {
			shown = shown
		}
		return
	}
	// 检查是否有人吃
	hand := m.playerTile[(currentOrder+1)%4].HandTiles
	if m.rules.CanChow(hand, msg.Tile) {

	}
	// 都没有
	handTile := m.playerTile[msg.TableOrder].HandTiles
	dropTile := m.playerTile[msg.TableOrder].DropTiles
	m.playerTile[msg.TableOrder].DropTiles = append(dropTile, msg.Tile)
	m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(handTile, msg.Tile, 1)
}

func (m *MahjongManager) handleChow(msg message.GameMsgRecv) {

}

// handlePong收到碰的消息，通知每个玩家，并等待该玩家发出打牌的消息
func (m *MahjongManager) handlePong(msg message.GameMsgRecv) {
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles
	// 更新明牌
	m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
		ShownType: model.Pong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile},
	})
	// 更新手牌
	newHand := model.RemoveTile(hand, msg.Tile, 2)
	m.playerTile[msg.TableOrder].HandTiles = newHand
	// 发送消息
	msgSend := message.GameMsgSend{
		MsgType:          message.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      m.currentTableOrder,
		CurrentTile:      nil, //其他玩家不会管， 当前玩家会忽略
		AvailableActions: []int{message.Discard},
		LastTurn:         m.currentTableOrder,
		LastAction:       message.Pong,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleWin(msg message.GameMsgRecv) {

}

func (m *MahjongManager) handleCancel(msg message.GameMsgRecv) {
	switch m.lastAction {
	case message.Win:
		//TODO: 检查碰、杠
	case message.Kong:
		fallthrough
	case message.Pong:
		//TODO: 检查吃
	case message.Chow:
		//TODO: 重新抓牌
	}
}

func (m *MahjongManager) getAvailableActions(hand []model.Tile, shown []model.ShownTile, newTile model.Tile) []int {
	availableActions := make([]int, 0)
	if m.rules.CanChow(hand, newTile) {
		availableActions = append(availableActions, message.Chow)
	}
	if m.rules.CanPong(hand, newTile) {
		availableActions = append(availableActions, message.Pong)
	}
	if m.rules.CanExposedKong(hand, newTile) {
		availableActions = append(availableActions, message.Kong)
	}
	if m.rules.CanWin(hand, shown, newTile) {
		availableActions = append(availableActions, message.Win)
	}
	return availableActions
}

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
