package block_query

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	. "github.com/EthereumHD/Scan/src/const"
	"github.com/EthereumHD/Scan/src/model"
	. "github.com/EthereumHD/Scan/src/model"
	"fmt"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
	log "qoobing.com/utillib.golang/log"
)

type InputHashReq struct {
	Hash string `json:"hash" form:"hash"`
}

type OutputHashRsp struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	model.BlockDetail
}

func Get_by_hash(cc echo.Context) error {
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Mysql()

	//Step 2. parameters initial

	rsp := OutputHashRsp{
		ErrNo:  0,
		ErrMsg: "success",
	}

	argc := new(InputHashReq)

	if err := c.BindInput(argc); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}
	log.Debugf("receive Get_by_hash: %+v", argc)

	//检查参数

	//查询区块

	//get transcation from chain
	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	chain_block, err := webthree.Eth.GetBlockByHash(argc.Hash, false)
	if err != nil {
		if err.Error() == EMPTY_RSP {
			log.Debugf("GetBlockByHash:%d from chain is NULL", argc.Hash)
			return c.RESULT_ERROR(BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(ERR_RPC_ERROR, err.Error())
	}

	poc, err := webthree.Eth.GetBlockPocByNumber(chain_block.Number)
	if err != nil {
		if err.Error() == EMPTY_RSP {
			log.Debugf("GetBlockPocByNumber:%d from chain is NULL", chain_block.Number.Int64())
			return c.RESULT_ERROR(BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(ERR_RPC_ERROR, err.Error())
	}

	//查询区块,数据库查询
	blocks, err := GetBlockByHash(c.Mysql(), argc.Hash)
	if err != nil {
		log.Debugf("GetBlockByHash error:%s,hash:%s", err.Error(), argc.Hash)
		return c.RESULT_ERROR(GET_BLOCKS_ERROR, fmt.Sprintf("GetBlockByHash error:%s,hash:%s", err.Error(), argc.Hash)) //c.RESULT(rsp)
	}
	//包装参数
	if len(blocks) == 0 {
		log.Debugf("the block not exist,hash:%s", argc.Hash)
		return c.RESULT_ERROR(BLOCK_OR_TRANS_NOT_EXIST, fmt.Sprintf("the block not exist,hash:%s", argc.Hash))
	}
	block := blocks[0]
	rsp.Height = block.F_block
	rsp.Hash = chain_block.Hash
	rsp.Transactions = block.F_txn
	rsp.Timestamp = chain_block.Timestamp.Int64()
	rsp.BlockReward = block.F_reward
	rsp.BlockFees = block.F_fees
	rsp.DeadLine = poc.Deadline.String()
	rsp.ExtraData = chain_block.ExtraData
	rsp.GasLimit = chain_block.GasLimit.String()
	rsp.GasUsed = chain_block.GasUsed.String()
	rsp.Miner = chain_block.Miner
	rsp.Nonce = chain_block.Nonce.String()
	rsp.ParentHash = chain_block.ParentHash
	rsp.Scoop = poc.ScoopNumber.String()
	rsp.Size = chain_block.Size.Int64()
	rsp.TotalDifficult = chain_block.TotalDifficult.String()
	rsp.Difficult = chain_block.Difficulty.String()

	//返回结果
	return c.RESULT(rsp)
}
