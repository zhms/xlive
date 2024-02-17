package model

import "github.com/shopspring/decimal"

type XLotteryOrder struct {
	Id         int64           `json:"id" gorm:"column:id"`
	SellerId   int             `json:"seller_id" gorm:"column:seller_id"`
	ChannelId  int             `json:"channel_id" gorm:"column:channel_id"`
	UserId     int             `json:"user_id" gorm:"column:user_id"`
	GameId     int             `json:"game_id" gorm:"column:game_id"`
	RoomLevel  int             `json:"room_level" gorm:"column:room_level"`
	Symbol     string          `json:"symbol" gorm:"column:symbol"`
	BetAmount  decimal.Decimal `json:"bet_amount" gorm:"column:bet_amount"`
	BetArea    string          `json:"bet_area" gorm:"column:bet_area"`
	BlockNum   int             `json:"block_num" gorm:"column:block_num"`
	State      int             `json:"state" gorm:"column:state"`
	TxId       string          `json:"tx_id" gorm:"column:tx_id"`
	Period     string          `json:"period" gorm:"column:period"`
	CreateTime string          `json:"create_time" gorm:"column:create_time"`
}
