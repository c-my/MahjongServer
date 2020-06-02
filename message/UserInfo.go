package message

type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
}
