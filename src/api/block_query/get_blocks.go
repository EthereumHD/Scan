package block_query

import (
	"fmt"
	"github.com/labstack/echo"
	. "github.com/EthereumHD/Scan/src/apicontext"
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/model"
	"qoobing.com/utillib.golang/log"
)

type InputReq struct {
	PageIndex int `json:"pageIndex" form:"pageIndex"` //范围起点
	PageSize  int `json:"pageSize" form:"pageSize"`   //范围重点
}

type OutputRsp struct {
	ErrNo  int       `json:"err_no"`
	ErrMsg string    `json:"err_msg"`
	Count  int64     `json:"count"` //区块个数
	Blocks BlockList `json:"blocks"`
}

type BlockInfo struct {
	BlockNumber int64  `json:"block_number"`
	Timestamp   int64  `json:"timestamp"`
	BlockMiner   string `json:"block_miner"`
	BlockReward string `json:"block_reward"`
	BlockFees   string `json:"block_fees"`
	GasUsed     string `json:"gas_used"`
	GasLimit    string `json:"gas_limit"`
}

type BlockList []BlockInfo

func Get_Blocks(cc echo.Context) error {
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
	//查询区块,数据库查询
	count, err := GetActiveBlockNum(c.Mysql())
	if err != nil {
		log.Debugf("GetActiveBlockNum error:", err.Error())
		return c.RESULT_ERROR(BLOCK_COUNT_ERROR, fmt.Sprintf("GetActiveBlockNum error:%s", err.Error())) //c.RESULT(rsp)
	}
	rsp.Count = count

	offset := (argc.PageIndex - 1) * argc.PageSize
	size := argc.PageSize
	blocks, err := GetRecentBlocks(c.Mysql(), offset, size)
	if err != nil {
		log.Debugf("GetRecentBlocks error:%s", err.Error())
		return c.RESULT_ERROR(GET_BLOCKS_ERROR, fmt.Sprintf("GetRecentBlocks error:%s", err.Error())) //c.RESULT(rsp)
	}
	//包装参数
	for _, block := range blocks {
		var blockInfo BlockInfo
		blockInfo.BlockNumber = block.F_block
		blockInfo.BlockMiner = block.F_miner
		blockInfo.BlockReward = block.F_reward
		blockInfo.Timestamp = block.F_timestamp
		blockInfo.BlockFees=block.F_fees
		blockInfo.GasUsed=block.F_gas_used
		blockInfo.GasLimit=block.F_gas_limit
		rsp.Blocks = append(rsp.Blocks, blockInfo)
	}
	//返回结果
	return c.RESULT(rsp)
}
