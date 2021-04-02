package components

import (
	"bytes"
	"github.com/clakeboy/golib/utils"
)

//命令 command
const (
	CMDPing     = 0x0001 //心跳
	CMDPong     = 0x0002 //心跳回复
	CMDAuth     = 0x0003 //认证登录
	CMDAuthCode = 0x0004 //认证结果
	CMDClose    = 0x00FF //关闭连接
)

//创建流式消息内容,加入内容长度
func BuildStreamData(data []byte) []byte {
	var buf bytes.Buffer
	lenByte := utils.IntToBytes(len(data), 32)
	buf.Write(lenByte)
	buf.Write(data)
	return buf.Bytes()
}

//解开内容列表
func ParseStreamData(data []byte) [][]byte {
	var list [][]byte
	idx := 0
	dataLength := len(data)
	for {
		if idx+4 > dataLength {
			break
		}
		length := utils.BytesToInt(data[idx : idx+4])
		idx += 4
		if idx+length > dataLength {
			break
		}
		list = append(list, data[idx:idx+length])
		idx += length
		if idx == dataLength {
			break
		}
	}
	return list
}
