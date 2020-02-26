/***********************************************************************
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
//******
// Filename:
// Description:
// Author:
// CreateTime:
/***********************************************************************/
package _const

const (
	ERR_DATABASE_ERROR        = -10001
	ERR_DATABASE_SELECT_ERROR = -10002
	ERR_DATABASE_SAVE_ERROR   = -10003

	ERR_REDIS_ERROR     = -10010
	ERR_REDIS_SET_ERROR = -10011
	ERR_REDIS_GET_ERROR = -10012
	ERR_REDIS_DEL_ERROR = -10013

	ERR_INNER_ERROR = -99999

	ERR_PARAMETER_INVALID = 10000

	ERR_RPC_ERROR = 20000

	BLOCK_COUNT_ERROR        = 30000
	GET_BLOCKS_ERROR         = 30001
	BLOCK_OR_TRANS_NOT_EXIST = 30002
	REPEAT_TRANSACTION       = 30400

	TRANSACTION_COUNT_ERROR = 40000
	GET_TRANSACTIONS_ERROR  = 40001
)

//status
const (
	ILLEGAL = iota
	NORMAL
	FORK
)

const (
	TCPKeepAliveTime               = 5 * 60 //5min
	MORTGAGECONTRACTADDR           = "0x0000000000000000000000000000000000000081"
	MORTGAGECONTRACT_FUNC_MORTGAGE = "0x43794dda"
	MORTGAGECONTRACT_FUNC_REDEEM   = "0x1e9a6950"
	ONEDAYBLOCK                    = 480
)

//
const (
	HTTPOK               = 200
	COOKIE_NAME_USERINFO = "USATK"
	DATA_NOT_EXIST       = "data not exist"
	EMPTY_RSP            = "Empty response"
)

//redis key
const (
	UBBEY_RATE = "UBBEY_RATE_KEY"
)
