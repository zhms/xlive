package service_hash

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"xclientapi/model"
	"xclientapi/server"
	"xcom/cuser"
	"xcom/edb"
	"xcom/enum"
	"xcom/utils"

	"github.com/beego/beego/logs"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type ServiceHash struct {
}

func (this *ServiceHash) Init() {
}

// 哈希下注
type HashBetReq struct {
	GameId    int             `json:"game_id"`    // 游戏Id
	RoomLevel int             `json:"room_level"` // 房间等级
	Symbol    string          `json:"symbol"`     // 币种
	Amount    decimal.Decimal `json:"amount"`     // 下注金额
	Area      string          `json:"area"`       // 下注区域
}

type HashBetRes struct {
	Amount decimal.Decimal `json:"amount"` // 剩余余额
	From   string          `json:"from"`   // 转出地址
	To     string          `json:"to"`     // 转入地址
	TxId   string          `json:"tx_id"`  // 交易Id
}

func (this *ServiceHash) HashBet(token *server.TokenData, host string, reqdata *HashBetReq) (response *HashBetRes, merr map[string]interface{}, err error) {
	reqdata.Symbol = "usdt"
	rediskey := fmt.Sprintf("%d_%d_%d_%d", token.SellerId, token.ChannelId, reqdata.GameId, reqdata.RoomLevel)
	strgameinfo, err := server.Redis(enum.RedisOrderRw).Client().HGet(context.Background(), "game:info", rediskey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetGameInfo error:", err.Error())
		return nil, nil, err
	}
	if errors.Is(err, redis.Nil) {
		return nil, enum.GameNotFound, nil
	}
	if strgameinfo == "" {
		return nil, enum.GameNotFound, nil
	}

	userdata := cuser.GetCacheData(token.UserId)
	if userdata == nil {
		return nil, enum.GetUserInfoError, nil
	}

	if userdata.Amount("usdt").LessThan(reqdata.Amount) {
		return nil, enum.AmountNotEnough, nil
	}
	posturl := "http://10.10.234.35:1524/api/send"
	resp, err := http.Post(posturl, "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		logs.Error("http post error1: %v", err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logs.Error("http post error2: %v", resp.StatusCode)
		return nil, enum.HashTransferError, nil
	}
	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		logs.Error("http post error3: %v", err)
		return nil, nil, err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(body.Bytes(), &data)
	if err != nil {
		logs.Error("http post error4: %v", err)
		return nil, nil, err
	}
	txid := utils.ToString(data["txid"])
	id, e := server.Redis(enum.RedisOrderRw).Client().LPop(context.Background(), "uniqueid").Result()
	if e != nil {
		logs.Error("get uniqueid error: %v", err)
		return nil, nil, err
	}
	result := cuser.AlterUserAmount(token.UserId, "usdt", reqdata.Amount.Neg(), decimal.Zero, decimal.Zero, enum.AmountChangeType_HashBet, id)
	if result == 4 {
		return nil, enum.AmountNotEnough, nil
	}
	if result != 0 {
		return nil, enum.BetError, nil
	}
	userdata = cuser.GetCacheData(token.UserId)
	if userdata == nil {
		return nil, enum.GetUserInfoError, nil
	}
	order := &model.XHashOrder{}
	order.Id = utils.ToInt64(id)
	order.SellerId = token.SellerId
	order.ChannelId = token.ChannelId
	order.UserId = token.UserId
	order.GameId = reqdata.GameId
	order.RoomLevel = reqdata.RoomLevel
	order.Symbol = reqdata.Symbol
	order.BetAmount = reqdata.Amount
	order.BetArea = reqdata.Area
	order.TxId = txid
	order.FromAddress = utils.ToString(data["from"])
	order.ToAddress = utils.ToString(data["to"])
	order.CreateTime = utils.Now()
	bytes, _ := json.Marshal(order)

	server.Redis(enum.RedisOrderRw).Client().RPush(context.Background(), fmt.Sprintf("order_yue:%v", order.TxId), string(bytes))
	server.Redis(enum.RedisOrderRw).Client().Expire(context.Background(), fmt.Sprintf("order_yue:%v", order.TxId), time.Duration(60*60*24*7)*time.Second)

	logs.Debug("hash_bet:", string(bytes))

	response = &HashBetRes{}
	response.From = utils.ToString(data["from"])
	response.To = utils.ToString(data["to"])
	response.TxId = utils.ToString(data["txid"])
	response.Amount = userdata.Amount("usdt")
	return response, nil, err
}

// 彩票下一期
type LotteryNextReq struct {
	GameId int `form:"game_id"` // 游戏Id
}

type LotteryNextRes struct {
	Period   string `json:"period"`    // 期号
	OpenTime string `json:"open_time"` // 开奖时间
	NowTime  string `json:"now_time"`  // 当前时间
}

func (this *ServiceHash) LotteryNext(reqdata *LotteryNextReq) (reponse *LotteryNextRes, merr map[string]interface{}, err error) {
	rediskey := fmt.Sprintf("period:%d", reqdata.GameId)
	strnext, err := server.Redis(enum.RedisOrderRw).Client().LIndex(context.Background(), rediskey, 0).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetGameInfo error:", err.Error())
		return nil, nil, err
	}
	if errors.Is(err, redis.Nil) {
		return nil, enum.GameNotFound, nil
	}
	if strnext == "" {
		return nil, enum.GameNotFound, nil
	}
	reponse = &LotteryNextRes{}
	err = json.Unmarshal([]byte(strnext), reponse)
	if err != nil {
		logs.Error("LotteryNext:", err.Error())
		return nil, nil, err
	}
	reponse.NowTime = utils.Now()
	return reponse, nil, err
}

// 彩票下注
type LotteryBetReq struct {
	GameId    int             `validate:"required" json:"game_id"`    // 游戏Id
	RoomLevel int             `validate:"required" json:"room_level"` // 房间等级
	Period    string          `validate:"required" json:"period"`     // 期号
	Symbol    string          `validate:"required" json:"symbol"`     // 币种
	Amount    decimal.Decimal `validate:"required" json:"amount"`     // 下注金额
	BetArea   string          `validate:"required" json:"bet_area"`   // 下注区域
}

type LotteryBetRes struct {
	Id     int64           `json:"id"`     // 订单Id
	Amount decimal.Decimal `json:"amount"` // 剩余金额
}

func (this *ServiceHash) LotteryBet(token *server.TokenData, host string, reqdata *LotteryBetReq) (reponse *LotteryBetRes, merr map[string]interface{}, err error) {
	reqdata.Symbol = edb.Usdt
	rediskey := fmt.Sprintf("period:%d", reqdata.GameId)
	strnext, err := server.Redis(enum.RedisOrderRw).Client().LIndex(context.Background(), rediskey, 0).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error("GetGameInfo error:", err.Error())
		return nil, nil, err
	}
	if errors.Is(err, redis.Nil) {
		return nil, enum.GameNotFound, nil
	}
	if strnext == "" {
		return nil, enum.GameNotFound, nil
	}
	next := &utils.XMap{RawData: map[string]interface{}{}}
	json.Unmarshal([]byte(strnext), &next.RawData)
	if reqdata.Period != next.String(edb.Period) {
		return nil, enum.GamePeriodError, nil
	}
	t := utils.LocalTimeToTimeStamp(next.String(edb.OpenTime))
	if time.Now().Unix() >= t-5 {
		return nil, enum.GamePeriodOver, nil
	}
	id, e := server.Redis(enum.RedisOrderRw).Client().LPop(context.Background(), "uniqueid").Result()
	if e != nil {
		logs.Error("get uniqueid error: %v", err)
		return nil, nil, err
	}
	result := cuser.AlterUserAmount(token.UserId, reqdata.Symbol, reqdata.Amount.Neg(), decimal.Zero, decimal.Zero, enum.AmountChangeType_LotteryBet, id)
	if result == 4 {
		return nil, enum.AmountNotEnough, nil
	}
	if result != 0 {
		return nil, enum.BetError, nil
	}
	order := &model.XLotteryOrder{}
	order.Id = utils.ToInt64(id)
	order.SellerId = token.SellerId
	order.ChannelId = token.ChannelId
	order.UserId = token.UserId
	order.GameId = reqdata.GameId
	order.RoomLevel = reqdata.RoomLevel
	order.Symbol = reqdata.Symbol
	order.BetAmount = reqdata.Amount
	order.BetArea = reqdata.BetArea
	order.State = enum.OrderState_WaitingOpen
	order.Period = reqdata.Period
	order.TxId = id
	order.CreateTime = utils.Now()
	ptime, err := time.Parse(utils.TimeLayout, order.CreateTime)
	if err != nil {
		logs.Error("LotteryBet time error:", err)
		return nil, nil, err
	}
	mtime := ptime.Format(utils.MounthLayout)
	db := server.Db(enum.DbOrderRW).Table(edb.TableOrderHash + mtime)
	db = db.Omit(edb.AuditTime, edb.RewardTime)
	err = db.Create(order).Error
	bytes, _ := json.Marshal(order)
	server.Redis(enum.RedisOrderRw).Client().RPush(context.Background(), fmt.Sprintf("order_lottery:%v:%v", reqdata.GameId, reqdata.Period), string(bytes))
	reponse = &LotteryBetRes{}
	reponse.Id = utils.ToInt64(id)
	reponse.Amount = cuser.GetCacheData(token.UserId).Amount(edb.Usdt)
	return reponse, nil, err
}

// 开奖历史
type LotteryHistoryReq struct {
}

type LotteryHistoryRes struct {
}

func (this *ServiceHash) LotteryHistory(reqdata *LotteryHistoryReq) (reponse *LotteryHistoryRes, merr map[string]interface{}, err error) {
	return reponse, nil, err
}
