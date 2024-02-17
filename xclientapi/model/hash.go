package model

import "github.com/shopspring/decimal"

type XHashOrder struct {
	Id          int64           `json:"id"`
	SellerId    int             `json:"seller_id"`
	ChannelId   int             `json:"channel_id"`
	UserId      int             `json:"user_id"`
	GameId      int             `json:"game_id"`
	RoomLevel   int             `json:"room_level"`
	Symbol      string          `json:"symbol"`
	BetAmount   decimal.Decimal `json:"bet_amount"`
	BetArea     string          `json:"bet_area"`
	TxId        string          `json:"tx_id"`
	FromAddress string          `json:"from_address"`
	ToAddress   string          `json:"to_address"`
	CreateTime  string          `json:"create_time"` // 创建时间
}
