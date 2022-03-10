package components

import (
	"bytes"
	"compress/gzip"
	"github.com/clakeboy/golib/utils"
	"io"
)

//命令 command
const (
	CMDPing     = 0x0001 //心跳
	CMDPong     = 0x0002 //心跳回复
	CMDAuth     = 0x0003 //认证登录
	CMDAuthCode = 0x0004 //认证结果
	CMDClose    = 0xFFFF //关闭连接
	CMDSysInfo  = 0x0005 //系统信息发送
	CMDShell    = 0x0006 //系统Shell执行
	CMDPty      = 0x0021 //打开终端机pty
	CMDPtyOpen  = 0x0022 //打开终端机pty
	CMDPtyClose = 0x0023 //打开终端机pty
	CMDPtyErr   = 0x0024 //打开终端机pty
	CMDFile     = 0x0010 //发送文件
	CMDDir      = 0x0011 //文件目录
)

// BuildStreamData 创建流式消息内容,加入内容长度
func BuildStreamData(data []byte) []byte {
	var buf bytes.Buffer
	lenByte := utils.IntToBytes(len(data), 32)
	buf.Write(lenByte)
	buf.Write(data)
	return buf.Bytes()
}

// ParseStreamData 解开内容列表
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

// Gzip 压缩数据
func Gzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	gbuf, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	defer gbuf.Close()
	_, err = gbuf.Write(data)
	if err != nil {
		return nil, err
	}
	err = gbuf.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnGzip 解压数据
func UnGzip(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	gbuf, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gbuf.Close()
	rBuf := bytes.NewBuffer([]byte{})
	tmp := make([]byte, 256)
	for {
		n, err := gbuf.Read(tmp)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, err
		}
		_, err = rBuf.Write(tmp[:n])
		if err != nil {
			return nil, err
		}
	}

	return rBuf.Bytes(), nil
}
