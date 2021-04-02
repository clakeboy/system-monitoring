package components

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/clakeboy/golib/utils"
)

var MainProtocol = []byte{0x03, 0xEE, 0xFE, 0x02}
var mask = []byte{0xEF, 0xEC}

//主通信信息流
type MainStream struct {
	Mask     []byte //掩码 2 byte
	Protocol []byte //通讯协议 4 byte
	Command  int    //命令 2 byte
	Content  []byte //内容
}

func NewMainStream() *MainStream {
	return &MainStream{
		Mask:     mask,
		Protocol: MainProtocol,
		Command:  0x0000,
		Content:  []byte{0x00},
	}
}

//生成通讯流
func (m *MainStream) Build() []byte {
	var stream bytes.Buffer
	stream.Write(m.Mask)
	stream.Write(m.Protocol)
	stream.Write(utils.IntToBytes(m.Command, 16))
	stream.Write(BuildStreamData(m.Content))
	stream.Write(m.Mask)
	return stream.Bytes()
}

//生成HEX通读流
func (m *MainStream) BuildHex() string {
	return hex.EncodeToString(m.Build())
}

//反解数据
func (m *MainStream) Parse(data []byte) error {
	if bytes.Compare(data[:2], mask) != 0 {
		return fmt.Errorf("invalid data")
	}
	if bytes.Compare(data[len(data)-2:], mask) != 0 {
		return fmt.Errorf("invalid data")
	}
	idx := 2
	m.Protocol = data[idx : idx+4]
	idx += 4
	cmd := data[idx : idx+2]
	m.Command = utils.BytesToInt(cmd)
	idx += 2
	contentLength := utils.BytesToInt(data[idx : idx+4])
	idx += 4
	m.Content = data[idx : idx+contentLength]

	return nil
}

//验证数据正确性
func (m *MainStream) Valid(data []byte) bool {
	protocol := data[2:6]
	return bytes.Compare(MainProtocol, protocol) == 0
}
