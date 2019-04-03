package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)
/*
将一个int64的整数：转为二进制后，每8bit一个byte.转为[]byte
 */
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	//将二进制数据写入
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	//转为[]byte并返回
	return buff.Bytes()
}


/*
Json字符串转为[] string数组
 */
func JsonToArray (jsonString string) []string {
	var sArr []string
	err := json.Unmarshal([]byte(jsonString), &sArr)
	if err != nil {
		log.Panic(err)
	}
	return sArr
}

//字节数组反转
func ReverseBytes(data []byte)  {
	for i,j := 0,len(data)-1;i <j ;i,j = i+1,j-1  {
		data[i],data[j] = data[j],data[i]
	}
}