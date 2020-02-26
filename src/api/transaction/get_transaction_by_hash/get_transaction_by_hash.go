package get_transaction_by_hash

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	"github.com/EthereumHD/Scan/src/const"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
	"math/big"
	"qoobing.com/utillib.golang/log"
)

type Input struct {
	Hash string `json:"hash" form:"hash" validate:"required"`
}

type TransactionDetail struct {
	TxHash           string `json:"tx_hash"`
	TxReceiptStatus  bool   `json:"tx_receipt_status"`
	Height           int64  `json:"height"`
	TimeStamp        int64  `json:"time_stamp"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasLimit         string `json:"gas_limit"`
	GasUsedByTx      string `json:"gas_used_by_tx"`
	GasPrice         string `json:"gas_price"`
	ActualTxCost     string `json:"actual_tx_cost"`
	Nonce            int64  `json:"nonce"`
	TransactionIndex int64  `json:"transaction_index"`
	InputData        string `json:"input_data"`
}

type Output struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	TransactionDetail
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
	if err := c.BindInput(&input); err != nil {
		return c.RESULT_PARAMETER_ERROR(err.Error())
	}

	//get transcation from chain
	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))

	transaction, err := webthree.Eth.GetTransactionByHash(input.Hash)
	if err != nil {
		if err.Error() == _const.EMPTY_RSP {
			log.Debugf("GetTransactionByHash:%s from chain is NULL", input.Hash)
			return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		//return c.RESULT_ERROR(_const.ERR_RPC_ERROR, err.Error()) //todo 这里pending 的transaction 也会返回
	}
	log.Debugf("GetTransactionByHash success,transcation%+v", transaction)

	blockbynumber, err := webthree.Eth.GetBlockByNumber(transaction.BlockNumber, false)
	if err != nil {
		if err.Error() == _const.EMPTY_RSP {
			log.Debugf("GetBlockByNumber:%d from chain is NULL", transaction.BlockNumber.Int64())
			return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(_const.ERR_RPC_ERROR, err.Error())
	}
	log.Debugf("GetBlockByNumber success")

	receipt, err := webthree.Eth.GetTransactionReceipt(input.Hash)
	if err != nil {
		if err.Error() == _const.EMPTY_RSP {
			log.Debugf("GetTransactionReceipt:%s from chain is NULL", input.Hash)
			return c.RESULT_ERROR(_const.BLOCK_OR_TRANS_NOT_EXIST, err.Error())
		}
		return c.RESULT_ERROR(_const.ERR_RPC_ERROR, err.Error())
	}
	log.Debugf("GetTransactionReceipt success")

	output.TxHash = input.Hash
	output.TxReceiptStatus = true
	output.Height = transaction.BlockNumber.Int64()
	output.TimeStamp = blockbynumber.Timestamp.Int64()
	output.From = transaction.From
	output.To = transaction.To
	output.Value = transaction.Value.String()
	output.GasLimit = transaction.Gas.String()
	output.GasUsedByTx = receipt.GasUsed.String()
	output.GasPrice = transaction.GasPrice.String()
	output.ActualTxCost = big.NewInt(0).Mul(receipt.GasUsed, transaction.GasPrice).String()
	output.Nonce = transaction.Nonce.Int64()
	output.InputData = transaction.Input
	output.TransactionIndex = transaction.TransactionIndex.Int64()

	// return
	output.ErrNo = 0
	output.ErrMsg = "success"

	return c.RESULT(output)
}
