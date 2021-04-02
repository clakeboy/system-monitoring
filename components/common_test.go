package components

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMainStream_BuildHex(t *testing.T) {
	var buf bytes.Buffer
	buf.Write(BuildStreamData([]byte("clake")))
	buf.Write(BuildStreamData([]byte("john")))
	buf.Write(BuildStreamData([]byte("lili")))
	data := buf.Bytes()
	fmt.Printf("%X\n", data)

	pData := ParseStreamData(data)

	for i, v := range pData {
		fmt.Printf("%d:%s\n", i, v)
	}
}
