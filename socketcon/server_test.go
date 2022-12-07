package socketcon

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	components2 "github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"system-monitoring/components"
	"testing"
)

func TestMultiData(t *testing.T) {
	str := "efec03eefe020022000000000100ecef3334efec03eefe020022000000000100ecef"
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
	components.CheckMultiStream(double[:])
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

	fmt.Printf("%c\n", 0x34)
	fmt.Printf("3344:%s", string([]byte{0x33, 0x34}))
}

func TestFuncValuePtr(t *testing.T) {
	data := utils.M{}
	client := utils.NewHttpClient()
	res, err := client.PostJsonString("https://twcins.cathaylife.cn/appuser/user/agentCompany/regist", data.ToJsonString())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}

var dirs = utils.M{
	"license_save_dir":      "/Volumes/tools/developer/train_data/car_license_front",
	"license_back_save_dir": "/Volumes/tools/developer/train_data/car_license_back",
	"id_front_save_dir":     "/Volumes/tools/developer/train_data/id_front",
	"id_back_save_dir":      "/Volumes/tools/developer/train_data/id_back",
}

func TestDownload(t *testing.T) {
	f, err := os.Open("/Users/clakeboy/Downloads/t_picc_brand.csv")
	if err != nil {
		t.Error(err)
		return
	}

	csvRed := csv.NewReader(f)
	rows, err := csvRed.ReadAll()
	if err != nil {
		t.Error(err)
		return
	}
	f.Close()
	pool := components2.NewPoll(8, func(obj ...interface{}) bool {
		fmt.Printf("begin download go index: %d \n", obj[1])
		row := obj[0].([]string)
		err = downloadFile(row[0], dirs["id_back_save_dir"].(string))
		if err != nil {
			fmt.Println("id_back_save_dir error", err, row[0])
		}
		err = downloadFile(row[1], dirs["id_front_save_dir"].(string))
		if err != nil {
			fmt.Println("id_front_save_dir error", err, row[1])
		}
		err = downloadFile(row[2], dirs["license_save_dir"].(string))
		if err != nil {
			fmt.Println("license_save_dir error", err, row[2])
		}
		err = downloadFile(row[3], dirs["license_back_save_dir"].(string))
		if err != nil {
			fmt.Println("license_back_save_dir error", err, row[3])
		}
		gp := len(obj[2].(*components2.GoroutinePool).Queue)
		fmt.Println("last files:", gp)
		return true
	})
	var list []interface{}

	for i, v := range rows {
		if i < 7991 {
			continue
		}
		list = append(list, v)
		//if len(list) >= 5 {
		//	break
		//}
	}

	pool.AddTaskInterface(list)
	pool.SetFinishCallback(func() {
		fmt.Println("download file successful!")
	})
	fmt.Println("start download files", len(list))
	pool.Start()

	//for i,row := range rows {
	//	fmt.Printf("begin download index: %d \n",i+1)
	//	err = downloadFile(row[0], dirs["id_back_save_dir"].(string))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = downloadFile(row[1], dirs["id_front_save_dir"].(string))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = downloadFile(row[2], dirs["license_save_dir"].(string))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	err = downloadFile(row[3], dirs["license_back_save_dir"].(string))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
}

var reg = regexp.MustCompile(`filename="(.+)"`)

func downloadFile(urlStr string, dir string) error {
	if !strings.ContainsAny(urlStr, "http://") {
		return nil
	}
	fmt.Println(path.Base(dir), urlStr)
	client := utils.NewHttpClient()
	res, err := client.Request("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("get file error :%v", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("request file status code error , status:%d", res.StatusCode)
	}

	name := path.Base(urlStr)
	savePath := fmt.Sprintf("%s/%s", dir, name)
	err = ioutil.WriteFile(savePath, res.Content, 0755)
	if err != nil {
		return err
	}
	fmt.Println("save file to", savePath)
	return nil
}

func TestCmdDir(t *testing.T) {
	dir := new(CMDDir)
	dir.Page = 1
	dir.Count = 100
	dir.Path = "/home"

	content := dir.Build()
	fmt.Printf("content:%X\n", content)

	newDir := new(CMDDir)
	err := newDir.Parse(content)
	if err != nil {
		t.Error(err)
	}
}
