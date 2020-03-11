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
		case message.ExposedKong:
		case message.ConcealedKong:
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
	newTile := m.wall.FrontDraw()
	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	model.SortTiles(m.playerTile[msg.TableOrder].HandTiles)
	//TODO:判断第一张牌能否碰、杠、胡

	var msgSend message.GameMsgSend

	msgSend = message.GameMsgSend{
		MsgType:          message.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: availableActions,
		LastTurn:         msg.TableOrder,
		LastAction:       message.Draw,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
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
}

func (m *MahjongManager) handleDiscard(msg message.GameMsgRecv) {
	//currentOrder := msg.TableOrder
	m.lastTableOrder = msg.TableOrder
	//从手牌中删除
	hand := m.playerTile[msg.TableOrder].HandTiles
	m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(hand, msg.Tile, 1)
	// 都没有
	// 加入打牌区
	dropTile := m.playerTile[msg.TableOrder].DropTiles
	m.playerTile[msg.TableOrder].DropTiles = append(dropTile, msg.Tile)
	//为下一位玩家发牌
	m.lastTableOrder = msg.TableOrder
	m.currentTableOrder = (msg.TableOrder + 1) % 4
	newTile := m.wall.FrontDraw()
	m.playerTile[m.currentTableOrder].HandTiles = append(m.playerTile[m.currentTableOrder].HandTiles, newTile)
	msgSend := message.GameMsgSend{
		MsgType:          message.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      (msg.TableOrder + 1) % 4,
		CurrentTile:      model.Tile{},
		AvailableActions: []int{message.Discard},
		LastTurn:         msg.TableOrder,
		LastAction:       message.Discard,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
	// end
	//// 检查是否有人胡牌
	//for i := 1; i <= 3; i++ {
	//	hand := m.playerTile[(currentOrder+i)%4].HandTiles
	//	shown := m.playerTile[(currentOrder+i)%4].ShownTiles
	//	if m.rules.CanWin(hand, shown, msg.Tile) {
	//		//send msg to this potential winner
	//		msgSend := message.GameMsgSend{
	//			MsgType:          -1,
	//			TableOrder:       -1,
	//			CurrentTurn:      (currentOrder + i) % 4,
	//			CurrentTile:      msg.Tile,
	//			AvailableActions: []int{message.Win},
	//			LastTurn:         msg.TableOrder,
	//			LastAction:       message.Discard,
	//			PlayerTile:       m.playerTile,
	//			WallCount:        m.wall.Length(),
	//		}
	//		m.gameSendCh <- msgSend
	//		return
	//	}
	//}
	//// 检查是否有人碰、杠
	//for i := 1; i <= 3; i++ {
	//	hand := m.playerTile[(currentOrder+i)%4].HandTiles
	//	if m.rules.CanExposedKong(hand, msg.Tile) {
	//		msgSend := message.GameMsgSend{
	//			MsgType:          message.GameMsgType,
	//			TableOrder:       -1,
	//			CurrentTurn:      (currentOrder + i) % 4,
	//			CurrentTile:      msg.Tile,
	//			AvailableActions: []int{message.Pong, message.ExposedKong},
	//			LastTurn:         msg.TableOrder,
	//			LastAction:       message.Discard,
	//			PlayerTile:       m.playerTile,
	//			WallCount:        m.wall.Length(),
	//		}
	//		m.gameSendCh <- msgSend
	//		return
	//	}
	//	if m.rules.CanPong(hand, msg.Tile) {
	//		msgSend := message.GameMsgSend{
	//			MsgType:          message.GameMsgType,
	//			TableOrder:       -1,
	//			CurrentTurn:      (currentOrder + i) % 4,
	//			CurrentTile:      msg.Tile,
	//			AvailableActions: []int{message.Pong},
	//			LastTurn:         msg.TableOrder,
	//			LastAction:       message.Discard,
	//			PlayerTile:       m.playerTile,
	//			WallCount:        m.wall.Length(),
	//		}
	//		m.gameSendCh <- msgSend
	//		return
	//	}
	//}
	//// 检查是否有人吃
	//nextHand := m.playerTile[(currentOrder+1)%4].HandTiles
	//if m.rules.CanChow(nextHand, msg.Tile) {
	//	msgSend := message.GameMsgSend{
	//		MsgType:          message.GameMsgType,
	//		TableOrder:       -1,
	//		CurrentTurn:      (currentOrder + 1) % 4,
	//		CurrentTile:      msg.Tile,
	//		AvailableActions: []int{message.Chow},
	//		LastTurn:         msg.TableOrder,
	//		LastAction:       message.Discard,
	//		PlayerTile:       m.playerTile,
	//		WallCount:        m.wall.Length(),
	//	}
	//	m.gameSendCh <- msgSend
	//	return
	//}
	//// 都没有
	//// 加入打牌区
	//dropTile := m.playerTile[msg.TableOrder].DropTiles
	//m.playerTile[msg.TableOrder].DropTiles = append(dropTile, msg.Tile)
	////为下一位玩家发牌
	//m.lastTableOrder = msg.TableOrder
	//m.currentTableOrder = (msg.TableOrder + 1) % 4
	//newTile := m.wall.FrontDraw()
	//m.playerTile[m.currentTableOrder].HandTiles = append(m.playerTile[m.currentTableOrder].HandTiles, newTile)
	//msgSend := message.GameMsgSend{
	//	MsgType:          message.GameMsgType,
	//	TableOrder:       -1,
	//	CurrentTurn:      (msg.TableOrder + 1) % 4,
	//	CurrentTile:      model.Tile{},
	//	AvailableActions: []int{message.Discard},
	//	LastTurn:         msg.TableOrder,
	//	LastAction:       message.Discard,
	//	PlayerTile:       m.playerTile,
	//	WallCount:        m.wall.Length(),
	//}
	//m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleChow(msg message.GameMsgRecv) {
	//更新手牌和明牌
	shown := m.playerTile[msg.TableOrder].ShownTiles
	// 删除手牌
	switch msg.ChowType {
	case message.LeftChow:
		mid := msg.Tile.GetRightTile()
		right := (*mid).GetRightTile()
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *mid, 1)
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *right, 1)
		m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
			ShownType: message.Chow,
			Tiles:     []model.Tile{msg.Tile, *mid, *right},
		})
	case message.MidChow:
		left := msg.Tile.GetLeftTile()
		right := msg.Tile.GetRightTile()
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *left, 1)
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *right, 1)
		m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
			ShownType: message.Chow,
			Tiles:     []model.Tile{*left, msg.Tile, *right},
		})
	case message.RightChow:
		mid := msg.Tile.GetLeftTile()
		left := (*mid).GetLeftTile()
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *mid, 1)
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *left, 1)
		m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
			ShownType: message.Chow,
			Tiles:     []model.Tile{*left, *mid, msg.Tile},
		})
	}
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:     message.GameMsgType,
		TableOrder:  -1,
		CurrentTurn: msg.TableOrder,
		//CurrentTile:      nil,
		AvailableActions: []int{message.Discard},
		LastTurn:         msg.TableOrder,
		LastAction:       message.Chow,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
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
		MsgType:     message.GameMsgType,
		TableOrder:  -1,
		CurrentTurn: m.currentTableOrder,
		//CurrentTile:      nil, //其他玩家不会管， 当前玩家会忽略
		AvailableActions: []int{message.Discard},
		LastTurn:         m.currentTableOrder,
		LastAction:       message.Pong,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleExposedKong(msg message.GameMsgRecv) {
	// 更新手牌和明牌
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles
	newHand := model.RemoveTile(hand, msg.Tile, 3)
	newShown := append(shown, model.ShownTile{
		ShownType: message.ExposedKong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile, msg.Tile},
	})
	m.playerTile[msg.TableOrder].HandTiles = newHand
	m.playerTile[msg.TableOrder].ShownTiles = newShown
	//为玩家发牌
	newTile := m.wall.BackDraw()
	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:          message.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: []int{message.Discard},
		LastTurn:         msg.TableOrder,
		LastAction:       message.ExposedKong,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleConcealedKong(msg message.GameMsgRecv) {
	// 更新手牌和明牌
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles
	newHand := model.RemoveTile(hand, msg.Tile, 3)
	newShown := append(shown, model.ShownTile{
		ShownType: message.ConcealedKong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile, msg.Tile},
	})
	m.playerTile[msg.TableOrder].HandTiles = newHand
	m.playerTile[msg.TableOrder].ShownTiles = newShown
	//为玩家发牌
	newTile := m.wall.BackDraw()
	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:          message.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: []int{message.Discard},
		LastTurn:         msg.TableOrder,
		LastAction:       message.ConcealedKong,
		PlayerTile:       m.playerTile,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleAddedKong(msg message.GameMsgRecv) {
	//更新手牌和明牌
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles

	m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(hand, msg.Tile, 1)
	for _, v := range shown {
		if v.ShownType == message.Pong && v.Tiles[0].Equals(msg.Tile) {
			v.ShownType = message.AddedKong
			v.Tiles = append(v.Tiles, msg.Tile)
			break
		}
	}

	//发牌
	newTile := m.wall.BackDraw()
	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:          message.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: []int{message.Discard},
		LastTurn:         msg.TableOrder,
		LastAction:       message.AddedKong,
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
	case message.ExposedKong:
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
		availableActions = append(availableActions, message.ExposedKong)
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
