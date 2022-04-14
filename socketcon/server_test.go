package socketcon

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestMultiData(t *testing.T) {
	str := "efec03eefe020022000000000100efec3334efec03eefe020022000000000100efec"
	//one,_ := hex.DecodeString(str)
	//two,_ := hex.DecodeString(str)
	//double := append(one,two...)
	double, _ := hex.DecodeString(str)
	//buf := bytes.NewReader(double)
	//buf := bytes.NewBuffer(double)
	//n,err := buf.ReadBytes(0xEC)
	//fmt.Println(hex.EncodeToString(n),err)
	//n,err = buf.ReadBytes(0xEC)
	//fmt.Println(hex.EncodeToString(n),err)
	//n,err = buf.ReadBytes(0xEC)
	//fmt.Println(hex.EncodeToString(n),err)
	//n,err = buf.ReadBytes(0xEC)
	//fmt.Println(hex.EncodeToString(n),err)
	//n,err = buf.ReadBytes(0xEC)
	//fmt.Println(hex.EncodeToString(n),err)
	var dataList [][]byte
	buf := bytes.NewBuffer([]byte{})
	read := bytes.NewBuffer(double)
	finish := false
	for {
		n, err := read.ReadBytes(0xec)
		if err != nil {
			break
		}
		fmt.Printf("%x ", n)
		fmt.Println(bytes.Equal(n[len(n)-2:], []byte{0xef, 0xec}))
		buf.Write(n)
		if finish && bytes.Equal(n[len(n)-2:], []byte{0xef, 0xec}) {
			fmt.Printf("finish: %x\n", buf.Bytes())
			dataList = append(dataList, buf.Bytes())
			buf = bytes.NewBuffer([]byte{})
			finish = false
			continue
		}

		if len(n) == 2 && !finish {
			finish = true
		}
	}
	fmt.Println("data list length:", len(dataList))
	for i, v := range dataList {
		fmt.Printf("%d:%x\n", i, v)
	}

	fmt.Printf("%d\n", 0x0022)
}

func TestFuncValuePtr(t *testing.T) {
	f := func() {}
	v := reflect.ValueOf(f)
	fmt.Println(v.Pointer())
	v = reflect.ValueOf(clakefunc)
	fmt.Println(v.Pointer())

}

//19562496
func clakefunc() error {
	return nil
}
