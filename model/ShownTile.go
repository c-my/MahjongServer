package model

const (
	Chow          = iota //吃
	Pong                 //碰
	ExposedKong          //明杠
	ConcealedKong        //暗杠
	AddedKong            //补杠
)

type ShownTile struct {
	ShownType int    `json:"shown_type"`
	Tiles     []Tile `json:"tiles"`
}

func GetChowCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == Chow {
			count++
		}
	}
	return count
}
func GetPongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == Pong {
			count++
		}
	}
	return count
}

func GetExposedKongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == ExposedKong {
			count++
		}
	}
	return count
}

func GetConcealedKongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == ConcealedKong {
			count++
		}
	}
	return count
}

func GetAddedKongCount(shownTiles []ShownTile) int {
	count := 0
	for _, v := range shownTiles {
		if v.ShownType == AddedKong {
			count++
		}
	}
	return count
}
