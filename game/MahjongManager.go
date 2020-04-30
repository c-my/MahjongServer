package game

import (
	"github.com/c-my/MahjongServer/config"
	"github.com/c-my/MahjongServer/message"
	"github.com/c-my/MahjongServer/model"
	"github.com/c-my/MahjongServer/rule"
	"log"
)

type MahjongManager struct {
	wall        *model.Wall
	playerTile  []model.PlayerTile
	rules       rule.MahjongRule
	firstPlayer int

	gameRecvCh   chan message.GameMsgRecv
	gameSendCh   chan message.GameMsgSend
	tableOrderCh chan int
	gameResultCh chan message.GameResultMsg
	getReadyCh   chan message.GetReadyMsg
	chatCh       chan message.ChatMsg
	exitCh       chan bool

	lastTableOrder    int
	currentTableOrder int

	cancelList  [4]bool
	readyList   [4]bool
	userList    []message.UserInfo
	lastMsgRecv message.GameMsgRecv
}

//NewMahjongManager return new Mahjong instance, with initiated game channels
// and tile-wall
func NewMahjongManager(gameRecvCh chan message.GameMsgRecv,
	gameSendCh chan message.GameMsgSend,
	tableOrderCh chan int,
	gameResultCh chan message.GameResultMsg,
	getReadyCh chan message.GetReadyMsg,
	chatCh chan message.ChatMsg,
	exitCh chan bool,
	r rule.MahjongRule) *MahjongManager {
	mahjong := MahjongManager{}
	mahjong.rules = r
	mahjong.wall = model.NewWall(r.GenerateTiles)
	mahjong.playerTile = make([]model.PlayerTile, 4)
	mahjong.firstPlayer = 0
	mahjong.gameRecvCh = gameRecvCh
	mahjong.gameSendCh = gameSendCh
	mahjong.tableOrderCh = tableOrderCh
	mahjong.gameResultCh = gameResultCh
	mahjong.getReadyCh = getReadyCh
	mahjong.chatCh = chatCh
	mahjong.exitCh = exitCh
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
		case e := <-m.exitCh:
			if e {
				return
			}
		case msg := <-gameRecvCh:
			log.Printf("player:[%d], action:[%d] ", msg.TableOrder, msg.Action)
			switch msg.Action {
			case config.Start:
				//inform each player their order
				m.tableOrderCh <- m.firstPlayer
				m.handleStart()
			case config.Ready:
				m.handleReady(msg)
			case config.Discard:
				m.handleDiscard(msg, false)
			case config.Chow:
				m.handleChow(msg)
			case config.Pong:
				m.handlePong(msg)
			case config.ExposedKong:
				m.handleExposedKong(msg)
			case config.ConcealedKong:
				m.handleConcealedKong(msg)
			case config.AddedKong:
				m.handleAddedKong(msg)
			case config.Win:
				m.handleWin(msg)
			case config.Cancel:
				m.handleCancel(msg)
			}

		}
	}
}

func (m *MahjongManager) handleStart() {
	m.resetCancelList()
	m.resetReadyList()

	//m.dealTile()
	m.dealTileTest()
	newTile := m.wall.FrontDraw()
	if newTile.IsEmpty() {
		m.gameResultCh <- m.getTieResult()
		return
	}
	//TODO:判断第一张牌能否碰、杠、胡
	availableActions := m.getAvailableActions(m.playerTile[m.firstPlayer].HandTiles, m.playerTile[m.firstPlayer].ShownTiles, newTile)
	m.playerTile[m.firstPlayer].HandTiles = append(m.playerTile[m.firstPlayer].HandTiles, newTile)
	model.SortTiles(m.playerTile[m.firstPlayer].HandTiles)
	var msgSend message.GameMsgSend

	msgSend = message.GameMsgSend{
		MsgType:          config.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      m.firstPlayer,
		CurrentTile:      newTile,
		AvailableActions: availableActions,
		LastTurn:         m.firstPlayer,
		LastAction:       config.Draw,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleReady(msg message.GameMsgRecv) {
	m.readyList[msg.TableOrder] = true
	msgSend := message.GetReadyMsg{
		MsgType:   config.GetReadyMsgType,
		ReadyList: m.readyList,
	}
	m.getReadyCh <- msgSend
	if canStartGame(m.readyList[:]) {
		//TODO: start game
		//inform each player their order
		m.tableOrderCh <- m.firstPlayer
		m.handleStart()
	}
}

func (m *MahjongManager) Draw(hand []model.Tile, shown []model.ShownTile, newTile model.Tile) {

}

func (m *MahjongManager) handleDiscard(msg message.GameMsgRecv, isJustCanceled bool) {
	currentOrder := msg.TableOrder
	m.lastTableOrder = msg.TableOrder
	lastAct := config.Cancel
	if !isJustCanceled { //从手牌中删除
		lastAct = config.Discard
		m.lastMsgRecv = msg
		hand := m.playerTile[msg.TableOrder].HandTiles
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(hand, msg.Tile, 1)
	}

	// 检查是否有人胡牌
	for i := 1; i <= 3; i++ {
		order := (currentOrder + i) % 4
		if m.cancelList[order] {
			continue
		}
		hand := m.playerTile[order].HandTiles
		shown := m.playerTile[order].ShownTiles
		if m.rules.CanWin(hand, shown, msg.Tile) {
			//send msg to this potential winner
			msgSend := message.GameMsgSend{
				MsgType:          config.GameMsgType,
				TableOrder:       -1,
				CurrentTurn:      order,
				CurrentTile:      msg.Tile,
				AvailableActions: []int{config.Win},
				LastTurn:         msg.TableOrder,
				LastAction:       lastAct,
				LastTile:         msg.Tile,
				PlayerTile:       m.playerTile,
				UserList:         m.userList,
				WallCount:        m.wall.Length(),
			}
			m.gameSendCh <- msgSend
			return
		}
	}
	// 检查是否有人碰、杠
	for i := 1; i <= 3; i++ {
		order := (currentOrder + i) % 4
		if m.cancelList[order] {
			continue
		}
		hand := m.playerTile[order].HandTiles
		if m.rules.CanExposedKong(hand, msg.Tile) {
			availableAction := make([]int, 0)
			chowTypes := make([]int, 0)
			availableAction = append(availableAction, config.Pong, config.ExposedKong)
			if i == 1 {
				canChow, types := m.rules.CanChow(hand, msg.Tile)
				if canChow {
					availableAction = append(availableAction, config.Chow)
					chowTypes = types
				}
			}
			msgSend := message.GameMsgSend{
				MsgType:          config.GameMsgType,
				TableOrder:       -1,
				CurrentTurn:      order,
				CurrentTile:      msg.Tile,
				AvailableActions: availableAction,
				LastTurn:         msg.TableOrder,
				LastAction:       lastAct,
				LastTile:         msg.Tile,
				ChowTypes:        chowTypes,
				PlayerTile:       m.playerTile,
				UserList:         m.userList,
				WallCount:        m.wall.Length(),
			}
			m.gameSendCh <- msgSend
			return
		}
		if m.rules.CanPong(hand, msg.Tile) {
			availableAction := make([]int, 0)
			chowTypes := make([]int, 0)
			availableAction = append(availableAction, config.Pong)
			if i == 1 {
				canChow, types := m.rules.CanChow(hand, msg.Tile)
				if canChow {
					availableAction = append(availableAction, config.Chow)
					chowTypes = types
				}
			}
			msgSend := message.GameMsgSend{
				MsgType:          config.GameMsgType,
				TableOrder:       -1,
				CurrentTurn:      order,
				CurrentTile:      msg.Tile,
				AvailableActions: availableAction,
				LastTurn:         msg.TableOrder,
				LastAction:       lastAct,
				LastTile:         msg.Tile,
				ChowTypes:        chowTypes,
				PlayerTile:       m.playerTile,
				UserList:         m.userList,
				WallCount:        m.wall.Length(),
			}
			m.gameSendCh <- msgSend
			return
		}
	}
	// 检查是否有人吃
	order := (currentOrder + 1) % 4

	nextHand := m.playerTile[order].HandTiles
	canChow, chowTypes := m.rules.CanChow(nextHand, msg.Tile)
	if !m.cancelList[order] && canChow {
		msgSend := message.GameMsgSend{
			MsgType:          config.GameMsgType,
			TableOrder:       -1,
			CurrentTurn:      order,
			CurrentTile:      msg.Tile,
			AvailableActions: []int{config.Chow},
			LastTurn:         msg.TableOrder,
			LastAction:       lastAct,
			LastTile:         msg.Tile,
			ChowTypes:        chowTypes,
			PlayerTile:       m.playerTile,
			UserList:         m.userList,
			WallCount:        m.wall.Length(),
		}
		m.gameSendCh <- msgSend
		return
	}
	// 都没有
	// 加入打牌区
	dropTile := m.playerTile[msg.TableOrder].DropTiles
	m.playerTile[msg.TableOrder].DropTiles = append(dropTile, msg.Tile)
	//为下一位玩家发牌
	m.lastTableOrder = msg.TableOrder
	m.currentTableOrder = (msg.TableOrder + 1) % 4
	newTile := m.wall.FrontDraw()
	if newTile.IsEmpty() {
		m.gameResultCh <- m.getTieResult()
		return
	}
	availableActions := m.getAvailableActions(m.playerTile[m.currentTableOrder].HandTiles, m.playerTile[m.currentTableOrder].ShownTiles, newTile)
	m.playerTile[m.currentTableOrder].HandTiles = append(m.playerTile[m.currentTableOrder].HandTiles, newTile)
	model.SortTiles(m.playerTile[m.currentTableOrder].HandTiles)
	msgSend := message.GameMsgSend{
		MsgType:          config.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      m.currentTableOrder,
		CurrentTile:      newTile,
		AvailableActions: availableActions,
		LastTurn:         msg.TableOrder,
		LastAction:       lastAct,
		LastTile:         msg.Tile,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
	// end
}

func (m *MahjongManager) handleChow(msg message.GameMsgRecv) {
	m.resetCancelList()
	//更新手牌和明牌
	shown := m.playerTile[msg.TableOrder].ShownTiles
	// 删除手牌
	switch msg.ChowType {
	case config.LeftChow:
		mid := msg.Tile.GetRightTile()
		right := (*mid).GetRightTile()
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *mid, 1)
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *right, 1)
		m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
			ShownType: config.Chow,
			Tiles:     []model.Tile{msg.Tile, *mid, *right},
		})
	case config.MidChow:
		left := msg.Tile.GetLeftTile()
		right := msg.Tile.GetRightTile()
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *left, 1)
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *right, 1)
		m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
			ShownType: config.Chow,
			Tiles:     []model.Tile{*left, msg.Tile, *right},
		})
	case config.RightChow:
		mid := msg.Tile.GetLeftTile()
		left := (*mid).GetLeftTile()
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *mid, 1)
		m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(m.playerTile[msg.TableOrder].HandTiles, *left, 1)
		m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
			ShownType: config.Chow,
			Tiles:     []model.Tile{*left, *mid, msg.Tile},
		})
	}
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:     config.GameMsgType,
		TableOrder:  -1,
		CurrentTurn: msg.TableOrder,
		//CurrentTile:      nil,
		AvailableActions: []int{config.Discard},
		LastTurn:         msg.TableOrder,
		LastAction:       config.Chow,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

// handlePong收到碰的消息，通知每个玩家，并等待该玩家发出打牌的消息
func (m *MahjongManager) handlePong(msg message.GameMsgRecv) {
	m.resetCancelList()
	m.lastTableOrder = m.currentTableOrder
	m.currentTableOrder = msg.TableOrder

	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles
	// 更新明牌
	m.playerTile[msg.TableOrder].ShownTiles = append(shown, model.ShownTile{
		ShownType: config.Pong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile},
	})
	// 更新手牌
	newHand := model.RemoveTile(hand, msg.Tile, 2)
	m.playerTile[msg.TableOrder].HandTiles = newHand
	// 发送消息
	msgSend := message.GameMsgSend{
		MsgType:     config.GameMsgType,
		TableOrder:  -1,
		CurrentTurn: m.currentTableOrder,
		//CurrentTile:      nil, //其他玩家不会管， 当前玩家会忽略
		AvailableActions: []int{config.Discard},
		LastTurn:         m.currentTableOrder,
		LastAction:       config.Pong,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleExposedKong(msg message.GameMsgRecv) {
	m.resetCancelList()

	m.lastTableOrder = m.currentTableOrder
	m.currentTableOrder = msg.TableOrder
	// 更新手牌和明牌
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles
	newHand := model.RemoveTile(hand, msg.Tile, 3)
	newShown := append(shown, model.ShownTile{
		ShownType: config.ExposedKong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile, msg.Tile},
	})
	m.playerTile[msg.TableOrder].HandTiles = newHand
	m.playerTile[msg.TableOrder].ShownTiles = newShown
	//为玩家发牌
	newTile := m.wall.BackDraw()
	if newTile.IsEmpty() {
		m.gameResultCh <- m.getTieResult()
		return
	}
	availableActions := m.getAvailableActions(m.playerTile[msg.TableOrder].HandTiles, m.playerTile[msg.TableOrder].ShownTiles, newTile)

	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	model.SortTiles(m.playerTile[msg.TableOrder].HandTiles)
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:          config.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: availableActions,
		LastTurn:         msg.TableOrder,
		LastAction:       config.ExposedKong,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleConcealedKong(msg message.GameMsgRecv) {
	m.resetCancelList()

	// 更新手牌和明牌
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles
	newHand := model.RemoveTile(hand, msg.Tile, 4)
	newShown := append(shown, model.ShownTile{
		ShownType: config.ConcealedKong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile, msg.Tile},
	})
	m.playerTile[msg.TableOrder].HandTiles = newHand
	m.playerTile[msg.TableOrder].ShownTiles = newShown
	//为玩家发牌
	newTile := m.wall.BackDraw()
	if newTile.IsEmpty() {
		m.gameResultCh <- m.getTieResult()
		return
	}
	availableActions := m.getAvailableActions(m.playerTile[msg.TableOrder].HandTiles, m.playerTile[msg.TableOrder].ShownTiles, newTile)
	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	model.SortTiles(m.playerTile[msg.TableOrder].HandTiles)
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:          config.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: availableActions,
		LastTurn:         msg.TableOrder,
		LastAction:       config.ConcealedKong,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleAddedKong(msg message.GameMsgRecv) {
	m.resetCancelList()

	//更新手牌和明牌
	hand := m.playerTile[msg.TableOrder].HandTiles
	shown := m.playerTile[msg.TableOrder].ShownTiles

	m.playerTile[msg.TableOrder].HandTiles = model.RemoveTile(hand, msg.Tile, 1)
	pos := 0
	for i, v := range shown {
		if v.ShownType == config.Pong && v.Tiles[0].Equals(msg.Tile) {
			pos = i
			break
		}
	}
	newShown := make([]model.ShownTile, 0)
	newShown = append(newShown, shown[:pos]...)
	newShown = append(newShown, model.ShownTile{
		ShownType: config.AddedKong,
		Tiles:     []model.Tile{msg.Tile, msg.Tile, msg.Tile, msg.Tile},
	})
	newShown = append(newShown, shown[pos+1:]...)
	m.playerTile[msg.TableOrder].ShownTiles = newShown

	//发牌
	newTile := m.wall.BackDraw()
	if newTile.IsEmpty() {
		m.gameResultCh <- m.getTieResult()
		return
	}
	availableActions := m.getAvailableActions(m.playerTile[msg.TableOrder].HandTiles, m.playerTile[msg.TableOrder].ShownTiles, newTile)
	m.playerTile[msg.TableOrder].HandTiles = append(m.playerTile[msg.TableOrder].HandTiles, newTile)
	model.SortTiles(m.playerTile[msg.TableOrder].HandTiles)
	//发送消息
	msgSend := message.GameMsgSend{
		MsgType:          config.GameMsgType,
		TableOrder:       -1,
		CurrentTurn:      msg.TableOrder,
		CurrentTile:      newTile,
		AvailableActions: availableActions,
		LastTurn:         msg.TableOrder,
		LastAction:       config.AddedKong,
		PlayerTile:       m.playerTile,
		UserList:         m.userList,
		WallCount:        m.wall.Length(),
	}
	m.gameSendCh <- msgSend
}

func (m *MahjongManager) handleWin(msg message.GameMsgRecv) {
	m.resetCancelList()
	msgSend := message.GameResultMsg{
		MsgType:    config.GameResultMsgType,
		Winner:     msg.TableOrder,
		FinalTile:  msg.Tile,
		PlayerTile: m.playerTile,
		UserList:   m.userList,
	}
	m.gameResultCh <- msgSend
	if msg.TableOrder != m.firstPlayer {
		m.firstPlayer++
		m.shiftUserList()
	}
}

func (m *MahjongManager) handleCancel(msg message.GameMsgRecv) {
	m.cancelList[msg.TableOrder] = true
	m.handleDiscard(m.lastMsgRecv, true)
}

func (m *MahjongManager) getAvailableActions(hand []model.Tile, shown []model.ShownTile, newTile model.Tile) []int {
	availableActions := make([]int, 0)
	availableActions = append(availableActions, config.Discard)
	if m.rules.CanConcealedKong(hand, newTile) {
		availableActions = append(availableActions, config.ConcealedKong)
	} else if m.rules.CanAddedKong(shown, newTile) {
		availableActions = append(availableActions, config.AddedKong)
	}
	if m.rules.CanWin(hand, shown, newTile) {
		availableActions = append(availableActions, config.Win)
	}
	return availableActions
}

func (m *MahjongManager) resetCancelList() {
	m.cancelList[0] = false
	m.cancelList[1] = false
	m.cancelList[2] = false
	m.cancelList[3] = false
}

func (m *MahjongManager) AddUserInfo(info message.UserInfo) {
	m.userList = append(m.userList, info)
}

func (m *MahjongManager) resetReadyList() {
	m.readyList[0] = false
	m.readyList[1] = false
	m.readyList[2] = false
	m.readyList[3] = false
}

func (m *MahjongManager) getTieResult() message.GameResultMsg {
	return message.GameResultMsg{
		MsgType:    config.GameResultMsgType,
		Winner:     -1,
		FinalTile:  model.Tile{},
		PlayerTile: m.playerTile,
		UserList:   m.userList,
	}
}

func canStartGame(readyList []bool) bool {
	for _, ready := range readyList {
		if !ready {
			return false
		}
	}
	return true
}

func (m *MahjongManager) shiftUserList() {
	tmp := m.userList[3]
	for i := 3; i > 0; i-- {
		m.userList[i] = m.userList[i-1]
	}
	m.userList[0] = tmp
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

func (m *MahjongManager) dealTileTest() {

	for i := 1; i <= 6; i++ {
		m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, model.Tile{Suit: model.Character, Number: i})
		m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, model.Tile{Suit: model.Character, Number: i}, model.Tile{Suit: model.Character, Number: i})
	}
	for i := 1; i <= 4; i++ {
		m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, model.Tile{Suit: model.Bamboo, Number: i})
		m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, model.Tile{Suit: model.Bamboo, Number: i}, model.Tile{Suit: model.Bamboo, Number: i}, model.Tile{Suit: model.Bamboo, Number: i})
	}
	m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, model.Tile{Suit: model.Dragon, Number: 1})
	m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, model.Tile{Suit: model.Dragon, Number: 1})
	m.playerTile[0].HandTiles = append(m.playerTile[0].HandTiles, model.Tile{Suit: model.Character, Number: 1})
	m.playerTile[1].HandTiles = append(m.playerTile[1].HandTiles, model.Tile{Suit: model.Dragon, Number: 1})
	m.playerTile[2].HandTiles = append(m.playerTile[2].HandTiles, model.Tile{Suit: model.Dragon, Number: 1})

	for i := 0; i < 13; i++ {
		m.playerTile[3].HandTiles = append(m.playerTile[3].HandTiles, m.wall.FrontDraw())
	}

	for i := 0; i <= 3; i++ {
		m.playerTile[i].SortHand()
	}
}
