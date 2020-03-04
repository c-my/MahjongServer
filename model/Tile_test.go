package model

import "testing"

func TestNewTile(t *testing.T) {
	tile := NewTile(Character, 3)
	if tile.Suit != Character || tile.Number != 3 {
		t.Errorf("wrong attribute value")
	}
}
