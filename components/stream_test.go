package components

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"testing"
)

func TestBinToBytes(t *testing.T) {
	bye, err := hex.DecodeString("FF2C2E1F")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bye)

	fmt.Println(0xff2e)

	protocol := []byte{0x03, 0xEE, 0xFE, 0x02}
	fmt.Println(bytes.Compare(MainProtocol, protocol))
	fmt.Println(utils.IntToBytes(CMDClose, 16))
}

func TestMainStream_Build(t *testing.T) {
	cmd := NewMainStream()
	cmd.Command = CMDClose
	cmdStream := cmd.Build()
	fmt.Println(hex.EncodeToString(cmdStream))

	deCmd := NewMainStream()
	err := deCmd.Parse(cmdStream)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedCode(t *testing.T) {

}
