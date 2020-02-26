/***********************************************************************
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
//******
// Filename:
// Description:
// Author:
// CreateTime:
/***********************************************************************/
package main

import (
	"github.com/EthereumHD/Scan/src/apicontext"
	"github.com/EthereumHD/Scan/src/config"
	"github.com/labstack/echo"
	"os"
	"os/exec"
	"github.com/EthereumHD/Scan/src/model"
	"github.com/EthereumHD/Scan/src/util"
	"qoobing.com/utillib.golang/gls"
	"qoobing.com/utillib.golang/log"

	"github.com/EthereumHD/Scan/src/api"
	"github.com/EthereumHD/Scan/src/statistics/stats"
	"github.com/EthereumHD/Scan/src/sync"
)

func main() {

	model.InitDatabase()

	filePath, _ := exec.LookPath(os.Args[0])
	log.Debugf("Program file: %s", filePath)

	go stats.Start()
	go sync.StartSyncLastBlock()
	//go sync.StartSyncRate()
	//go sync.CheckReward()

	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		//timer.CountMining()
		return func(c echo.Context) error {
			cc := apicontext.New(c)
			req := cc.Request()
			cc.RecordTime()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = util.GetRandomString(12)
			}
			gls.SetGlsValue("logid", id)

			log.Debugf("apicontext created")

			return h(cc)
		}
	})

	////user
	//e.POST("/user/login", api.Login)
	//e.POST("/user/logout", api.LogOut)

	//block
	e.POST("/block/get_block_number", api.GetBlockNumber)
	e.POST("/block/get_by_height", api.GetBlockByHeight)
	e.POST("/block/get_blocks", api.GetBlocks)
	e.POST("/block/get_by_hash", api.GetBlockByHash)

	//transaction
	e.POST("/transaction/get_by_hash", api.GetTransactionByHash)
	e.POST("/transaction/get_transactions", api.GetTransactions)
	e.POST("/transaction/get_by_addr", api.GetByAddr)
	e.POST("/transaction/get_by_addr_and_type", api.GetByAddrAndType)
	e.POST("/transaction/get_by_height", api.GetByHeight)
	e.POST("/transaction/get_addr_pending", api.GetAddrPending)
	e.POST("/transaction/get_hash_pending", api.GetHashPending)

	//mining
	e.POST("/mining/get_mined_block_by_addr", api.GetMinedBlocks)
	e.POST("/mining/get_addr_mining_rewards", api.GetAddrMiningRewards)
	e.POST("/mining/get_mined_block_by_addr_and_date", api.GetMinedblockByAddrAndDate)

	//poc
	e.POST("/poc/get_exchange_rate", api.GetExchangeRate)
	e.GET("/poc/get_exchange_rate", api.GetExchangeRate)
	e.POST("/poc/get_summary", api.GetSummary)
	e.GET("/poc/get_summary", api.GetSummary)
	e.POST("/poc/get_balance", api.GetBalance)

	e.Logger.Fatal(e.Start(":" + config.Config().Port))

}
