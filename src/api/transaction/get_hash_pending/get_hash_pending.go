package get_hash_pending

import (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
	"qoobing.com/utillib.golang/log"
	"strings"
)

type Input struct {
	Hash string `json:"hash" form:"hash" validate:"required"`
}

type Output struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	TransactionDetail
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

	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	content, err := webthree.Txpool.Content()
	if err != nil {
		log.Fatalf("Content error:%s", err.Error())
		return err
	}

	for _, txmap := range content.Pending {
		for _, tx := range txmap {
			if tx.Hash == strings.ToLower(input.Hash) {
				output.TxHash = input.Hash
				output.From = tx.From
				output.To = tx.To
				output.Value = tx.Value.String()
				output.Nonce = tx.Nonce.Int64()
			}
		}
	}

	output.ErrNo = 0
	output.ErrMsg = "success"

	return c.RESULT(output)
}
