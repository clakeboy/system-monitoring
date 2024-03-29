package components

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/clakeboy/golib/utils"
)

var MainProtocol = []byte{0x03, 0xEE, 0xFE, 0x02}
var Mask = []byte{0xEF, 0xEC}

// MainStream 主通信信息流
type MainStream struct {
	Mask     []byte //掩码 2 byte
	EndMask  []byte //结束掩码
	Protocol []byte //通讯协议 4 byte
	Command  int    //命令 2 byte
	Content  []byte //内容
	gzip     bool   //是否gzip压缩
}

func NewMainStream() *MainStream {
	endMask := make([]byte, 2)
	endMask[0], endMask[1] = Mask[1], Mask[0]
	return &MainStream{
		Mask:     Mask,
		EndMask:  endMask,
		Protocol: MainProtocol,
		Command:  0x0000,
		Content:  []byte{0x00},
	}
}

// Gzip 设置是否使用gzip压缩
func (m *MainStream) Gzip(flag bool) {
	m.gzip = flag
}

// Build 生成通讯流
func (m *MainStream) Build() []byte {
	var stream bytes.Buffer
	stream.Write(m.Mask)
	stream.Write(m.Protocol)
	stream.Write(utils.IntToBytes(m.Command, 16))
	if m.gzip {
		stream.Write([]byte{0x01})
		zipData, err := Gzip(m.Content)
		if err != nil {
			return nil
		}
		contentLength := utils.IntToBytes(len(zipData), 32)
		stream.Write(contentLength)
		stream.Write(zipData)
	} else {
		stream.Write([]byte{0x00})
		contentLength := utils.IntToBytes(len(m.Content), 32)
		stream.Write(contentLength)
		stream.Write(m.Content)
	}
	stream.Write(m.EndMask)
	return stream.Bytes()
}

// BuildHex 生成HEX通读流
func (m *MainStream) BuildHex() string {
	return hex.EncodeToString(m.Build())
}

// Parse 反解数据
func (m *MainStream) Parse(data []byte) error {
	if !bytes.Equal(data[:2], m.Mask) {
		return fmt.Errorf("invalid data")
	}
	if !bytes.Equal(data[len(data)-2:], m.EndMask) {
		return fmt.Errorf("invalid data")
	}
	idx := 2
	m.Protocol = data[idx : idx+4]
	idx += 4
	cmd := data[idx : idx+2]
	m.Command = utils.BytesToInt(cmd)
	idx += 2
	gzip := data[idx : idx+1]
	m.gzip = gzip[0] == 0x01
	idx += 1
	contentLength := utils.BytesToInt(data[idx : idx+4])
	idx += 4
	if m.gzip {
		zipData, err := UnGzip(data[idx : idx+contentLength])
		if err != nil {
			return err
		}
		m.Content = zipData
	} else {
		m.Content = data[idx : idx+contentLength]
	}
	return nil
}

// Valid 验证数据正确性
func (m *MainStream) Valid(data []byte) bool {
	protocol := data[2:6]
	return bytes.Equal(MainProtocol, protocol)
}

func CheckMultiStream(data []byte) ([]*MainStream, error) {
	defer func() {
		if err := recover(); err != nil {
			// fmt.Println("check multi stream error:", err)
			fmt.Printf("check multi stream error:%v\n%X\n", err, data)
		}
	}()
	var dataList []*MainStream
	endMask := make([]byte, 2)
	endMask[0], endMask[1] = Mask[1], Mask[0]
	for {
		beginIdx := bytes.Index(data, Mask)
		if beginIdx == -1 {
			break
		}
		if beginIdx > 0 {
			data = data[beginIdx:]
		}
		endIdx := bytes.Index(data, endMask)
		if endIdx == -1 {
			break
		}
		msg := data[:endIdx+2]
		data = data[endIdx+2:]
		stream := NewMainStream()
		err := stream.Parse(msg)
		if err != nil {
			return nil, fmt.Errorf("check data error: %v, org: %x", err, msg)
		}
		dataList = append(dataList, stream)
	}
	return dataList, nil
}
