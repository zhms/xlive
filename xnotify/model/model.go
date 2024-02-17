package model

type BlockInfo struct {
	BlockNum  int64  `json:"block_num"`
	BlockTime string `json:"block_time"`
	BlockHash string `json:"block_hash"`
}
