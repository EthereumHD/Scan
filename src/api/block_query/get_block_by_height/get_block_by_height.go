package get_block_by_height

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	"github.com/EthereumHD/Scan/src/const"
	"github.com/EthereumHD/Scan/src/model"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
	"math/big"
	"qoobing.com/utillib.golang/log"
)

type Input struct {
	Height int64 `json:"height" form:"height" `
}

type Output struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	model.BlockDetail
}

func Main(cc echo.Context) error {

	//Step 1. init x
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Redis()
	c.Mysql()

	//Step 2. parameters initial
	var (
		input  Input
		output Output
	)
	output.ErrNo = 0
	output.ErrMsg = "success"

	if err := c.BindInput(&input); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}

	//get transcation from chain
	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	chain_block, err := webthree.Eth.GetBlockByNumber(big.NewInt(input.Height), false)
	if err != nil {
		if err.Error() == _const.EMPTY_RSP {
			log.Debugf("GetBlockByNumber:%d from chain is NULL", input.Height)
			return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(_const.ERR_RPC_ERROR, err.Error())
	}

	poc, err := webthree.Eth.GetBlockPocByNumber(big.NewInt(input.Height))
	if err != nil {
		if err.Error() == _const.EMPTY_RSP {
			log.Debugf("GetBlockPocByNumber:%d from chain is NULL", input.Height)
			return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(_const.ERR_RPC_ERROR, err.Error())
	}

	databases_block, err := (&model.Block{}).FindBlockByHeight(c.Mysql(), input.Height)
	if err != nil {
		if err.Error() == _const.DATA_NOT_EXIST {
			log.Debugf("FindBlockByHeight:%d from databases is NULL", input.Height)
			return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}

		return c.RESULT_ERROR(_const.ERR_DATABASE_ERROR, err.Error())
	}

	output.BlockDetail = model.BlockDetail{
		Height:         input.Height,
		Timestamp:      chain_block.Timestamp.Int64(),
		Transactions:   int64(len(chain_block.Transactions)),
		Hash:           chain_block.Hash,
		ParentHash:     chain_block.ParentHash,
		Miner:          chain_block.Miner,
		Difficult:      chain_block.Difficulty.String(),
		TotalDifficult: chain_block.TotalDifficult.String(),
		Size:           chain_block.Size.Int64(),
		GasUsed:        chain_block.GasUsed.String(),
		GasLimit:       chain_block.GasLimit.String(),
		Nonce:          chain_block.Nonce.String(),
		BlockReward:    databases_block.F_reward,
		BlockFees:      databases_block.F_fees,
		ExtraData:      chain_block.ExtraData,
		DeadLine:       poc.Deadline.String(),
		Scoop:          poc.ScoopNumber.String(),
	}
	//todo extradat scopp check

	return c.RESULT(output)
}
