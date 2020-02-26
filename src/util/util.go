/***********************************************************************
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
//******
// Filename:
// Description:
// Author:
// CreateTime:
/***********************************************************************/
package util

import (
	"math/rand"
	"time"
	"crypto/md5"
	"encoding/hex"
)

//////////////global assert//////////////////begin/////////////////////
func ASSERT(expr bool, errstr string) {
	if !expr {
		panic(errstr)
	}
}

//////////////global function//////////////////begin/////////////////////
func GetRandomString(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomCharacter(l int) string {
	str1 := "0123456789"
	str2 := "abcdefghijKlmnopqrstuvwxyz"
	str3 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str1 + str2 + str3)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func CheckDate(date string) bool {
	withNanos := "2006-01-02 15:04:05"
	_, err := time.ParseInLocation(withNanos, date, time.Local)
	if err != nil {
		return false
	} else {
		return true
	}
}

func MD5(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	encode := hex.EncodeToString(md5Ctx.Sum(nil))
	return encode
}

func init() {
}
