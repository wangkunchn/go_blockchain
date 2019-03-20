package main

import "os"

const DB_NAME = "blockchain.db"	//数据库名
const BLOCK_TABLE_NAME = "blocks" //表名
const LAST_BLOCK_HASH = "last"	//blockchain 中最后一个block  对应的key

func dbExists() bool {
	_, err:= os.Stat(DB_NAME)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
