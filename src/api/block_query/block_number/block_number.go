package block_number

import                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 (
	. "github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	"github.com/EthereumHD/Scan/src/const"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
)

type Output struct {
	ErrNo       int    `json:"err_no"`
	ErrMsg      string `json:"err_msg"`
	BlockNumber int64  `json:"block_number"`
}

func Main(cc echo.Context) error {

	//Step 1. init x
	c := cc.(ApiContext)
	defer c.PANIC_RECOVER()
	c.Redis()
	c.Mysql()

	//Step 2. parameters initial
	var (
		output Output
	)

	//get transcation from chain
	webthree := web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	number, err := webthree.Eth.GetBlockNumber()
	if err != nil {
		return c.RESULT_ERROR(_const.ERR_RPC_ERROR, err.Error())
	}

	output.ErrNo = 0
	output.ErrMsg = "success"
	output.BlockNumber = number.Int64()

	return c.RESULT(output)
}
