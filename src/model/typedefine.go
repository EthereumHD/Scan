/***********************************************************************
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
//******
// Filename:
// Description:
// Author:
// CreateTime:
/***********************************************************************/
package model

type BlockDetail struct {
	Height         int64  `json:"height"`
	Timestamp      int64  `json:"timestamp"`
	Transactions   int64  `json:"transactions"`
	Hash           string `json:"hash"`
	ParentHash     string `json:"parent_hash"`
	Miner          string `json:"miner"`
	Difficult      string `json:"difficult"`
	TotalDifficult string `json:"total_difficult"`
	Size           int64  `json:"size"`
	GasUsed        string `json:"gas_used"`
	GasLimit       string `json:"gas_limit"`
	Nonce          string `json:"nonce"`
	BlockReward    string `json:"block_reward"`
	BlockFees      string `json:"block_fees"`
	ExtraData      string `json:"extra_data"`
	DeadLine       string `json:"deadline"`
	Scoop          string `json:"scoop"`
}

type UbbeyRate struct {
	Eth float64 `json:"eth"`
	Btc float64 `json:"btc"`
	USD float64 `json:"usd"`
	KWR float64 `json:"kwr"`
}
