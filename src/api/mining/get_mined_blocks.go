package mining

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/model"
	"fmt"
	"github.com/labstack/echo"
	"qoobing.com/utillib.golang/log"
)

type InputReq struct {
	Addr      string `json:"addr" form:"addr"`
	PageIndex int    `json:"pageIndex" form:"pageIndex"` //范围起点
	PageSize  int    `json:"pageSize" form:"pageSize"`   //范围重点
}

type OutputRsp struct {
	ErrNo  int       `json:"err_no"`
	ErrMsg string    `json:"err_msg"`
	Count  int64     `json:"count"` //区块个数
	Blocks BlockList `json:"blocks"`
}

type BlockInfo struct {
	Hash        string `json:"hash"`
	BlockNumber int64  `json:"block_number"`
	Timestamp   int64  `json:"timestamp"`
	Txn         int64  `json:"txn"`
	//Difficulty  string `json:"difficulty"`
	BlockFees   string `json:"block_fees"`
	BlockReward string `json:"block_reward"`
	GasUsed     string `json:"gas_used"`
}

type BlockList []BlockInfo

func Get_mined_block_by_addr(cc echo.Context) error {
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Mysql()

	//Step 2. parameters initial

	rsp := OutputRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	argc := new(InputReq)

	if err := c.BindInput(argc); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}
	log.Debugf("receive Get_Blocks: %+v", argc)

	//检查参数
	if argc.PageIndex < 1 || argc.PageSize <= 0 {
		log.Debugf("param error")
		return c.RESULT_ERROR(ERR_PARAMETER_INVALID, "param error")
	}
	//TODO检查地址正确性

	//查询区块,数据库查询
	//sql := "F_miner = '" + argc.Addr + "" + "' "
	count, err := GetActiveBlockNumByAddr(c.Mysql(), argc.Addr)
	if err != nil {
		log.Debugf("GetActiveBlockNum error:%s,addr:%s", err.Error(), argc.Addr)
		return c.RESULT_ERROR(BLOCK_COUNT_ERROR, fmt.Sprintf("GetActiveBlockNum error:%s,addr:%s", err.Error(), argc.Addr)) //c.RESULT(rsp)
	}
	rsp.Count = count

	offset := (argc.PageIndex - 1) * argc.PageSize
	size := argc.PageSize
	blocks, err := GetBlocksByMinerAddr(c.Mysql(), argc.Addr, offset, size)
	if err != nil {
		log.Debugf("GetBlocksByMinerAddr error:%s,addr:%s", err.Error(), argc.Addr)
		return c.RESULT_ERROR(GET_BLOCKS_ERROR, fmt.Sprintf("GetBlocksByMinerAddr error:%s,addr:%s", err.Error(), argc.Addr)) //c.RESULT(rsp)
	}
	//包装参数
	for _, block := range blocks {
		var blockInfo BlockInfo
		blockInfo.Hash = block.F_hash
		blockInfo.BlockNumber = block.F_block
		blockInfo.Timestamp = block.F_timestamp
		blockInfo.Txn = block.F_txn
		blockInfo.BlockFees = block.F_fees
		blockInfo.BlockReward = block.F_reward
		blockInfo.GasUsed = block.F_gas_used

		rsp.Blocks = append(rsp.Blocks, blockInfo)
	}
	//返回结果
	return c.RESULT(rsp)
}
