package get_balance

import (
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"

	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	. "github.com/EthereumHD/Scan/src/const"
	"github.com/EthereumHD/Scan/src/model"
	"fmt"
	"go-web3/eth/block"
	"qoobing.com/utillib.golang/log"
)

type Input struct {
	Addr string `json:"addr" form:"addr" validate:"required"`
}

type Output struct {
	ErrNo        int    `json:"err_no"`
	ErrMsg       string `json:"err_msg"`
	Balance      string `json:"balance"`
	Transactions int64  `json:"transactions"`
	MinedBlocks  int64  `json:"mined_blocks"`
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

	//get balance from chain
	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	bal, err := webthree.Eth.GetBalance(input.Addr, block.LATEST)
	if err != nil {
		return c.RESULT_ERROR(ERR_RPC_ERROR, err.Error())
	}

	count, err := model.GetActiveBlockNumByAddr(c.Mysql(), input.Addr)
	if err != nil {
		log.Debugf("GetActiveBlockNumByAddr error:%s,addr:%s", err.Error(), input.Addr)
		return c.RESULT_ERROR(BLOCK_COUNT_ERROR, fmt.Sprintf("GetActiveBlockNumByAddr error:%s,addr:%s", err.Error(), input.Addr)) //c.RESULT(output)
	}
	output.MinedBlocks = count

	count, err = model.GetTransactionsCountByAddr(c.Mysql(), input.Addr)
	if err != nil {
		log.Debugf("GetTransactionsCountByAddr error:%s,addr:%s", err.Error(), input.Addr)
		return c.RESULT_ERROR(TRANSACTION_COUNT_ERROR, fmt.Sprintf("GetTransactionsCountByAddr error:%s,addr:%s", err.Error(), input.Addr)) //c.RESULT(output)
	}

	output.Transactions = count
	output.ErrNo = 0
	output.ErrMsg = "success"
	output.Balance = bal.String()

	return c.RESULT(output)
}
