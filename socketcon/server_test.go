package socketcon

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"system-monitoring/components"
	"testing"

	components2 "github.com/clakeboy/golib/components"
	"github.com/clakeboy/golib/utils"
)

func TestMultiData(t *testing.T) {
	str := "DDDE96448FA64A7B87F50FB9FEEDCD83C202611EF3A768EB06E253F47F2F5C44400ABE9836FF9B358A2B21500C1C51D846D8B4D1EC21F5F6E21331A4B6E2F7C6FA0F2432636298A3672B0818ECC235DD6EF865FB6BDF3DF47918A6A668BA6310B106B5663C641BEFB9FEE55FAE045594C8F2366B8C40E186DC3755FBCFFD6699E7E1D656CBDCDAE36B6131F009264AE18F30C13526894AF82D4CA45EC3F730217FC48486E00B4C21F84F306188728909FE02A7E53526460EC1707CE0049042D02B3D09244840E9335218C57B0EF20E55520265E43435D9252A49115505503CBD63858A497CD0F7AC0245D5E922FE0416B3A829F83D2D8C9880285E8A6AF95760C167B410F18F6849BA92558C90D006CE776811317A7B197C8316A42B5A41BCC610BEA045F6E528D7B42085E8035ED2FA4A5BAF855BF5B91A73FDA3AE46FBB879FD3F000000FFFFECEF3334EFEC03EEFE02000500000004CA1F8B08000000000002FF94964B6FEBCA0D80FF0BD782337CCC90A3650B14B89BF6E2F6165D04C1816C4D6221B26548727CD220FFBD1859721E27E7E6249BC824878F6F484A4FD0D450BA027669D7F58F503EC1D88D550B2547F58E4C3914503D544D5BADDB0425B1A10A392BE038A41A4A16F514D49B9C25BFA77E93F62394A62B2115638AE899B080DB3E2528D199794F31FBDD8CCD438292392A297A4F0534FB458AD1C491132BE0D4F4E99C675B1DF7754ED415B03EDEDEA67E8052D89361A40236D5663B67C5AA214801A7BE19D3BADADC4F67EAA61F1FA18C0E915EE9C6DD61520FDB6A8A84211BE4E3435BADA14424120A66050C7DDAB455B33BF35054275E4301C3713F6BA014E7D9A60887EA2E8DD974805245396A960EA7EAB0A4EA0AD874BB5D33B6CDAE19A7B2D5A370B04531A6BA1AA08C91A3573373056C9BBBED7C51F3AF33DC4CA83BBD68DAEE7451E4A02F9AFCEBA2DA558743CE45CC05925CF6C3AE6ADB6EB3348347135672A696AF7ED69E3B0085CD63BEBB59BCD91EF7F7CB2131F51C296779BC4B19C7F02AF14574C9E42269FE97BBCD45454FCF05ECD3F8ADD98FA9BFAD3609CAEB2768F675FA0E2516B01B8F5006EF73ABEEAB5D8212DA0E0AD8567D7DAAFA54D5750F254001B76D753740790DC7036438DD616A8C9B02B2CD30F99DAD9174E5566E855706CF37CFC52520CD01D13B778997C6ADFB31A27325869253E942295AD6E18714D67D57D59B6A18A180DDB11D9BE9F9C37C9456682B545979BC229793BA99B97493E19C4AD5B6D9F1E398867F4F8318D53CAB620838CBFF489B07282347440B48D1729F6EEED3389FB0182385807891CF271C466FE6A880D4F7CD7EBAB1D4F7DD713C0F57DF1D66697E5CC4B7CD6D378BF3E3599C93DF1C8EDF8E43CA5B27675DD22A7BF726CC81D5F2E4B7CD304279CD2B17498D3DC5E81C062968E58CA73F65F59EE9E6394FF7707FF678FD04876ADC42095719FA303E1EA67BFA3ED255FA3E3214CBA21342711263CC0377EE432666179C789165CD69C883E11DBD5F731857E8A318A1283911CC2BACABD3F0E7D93D054211B748FF33798BE24C17D13FA698E449D9F36BBB4B0C5E798BC84218D8998F92DB7129EF743A7D522163888C96F7072D2552608FAA82CE87CB2A578AC84EDDB2DD5F8AA41547758E5D349430ED803755620CC1997B57A6788A8EDED689919C5AB40F0BA5153BF66A214A0828DE7C6E93E95673933FC1435DE57F79C2FEDE1DF391E028FA0C7397FABB54FFF1A221F58A3C2FF9C53ACF82F8C5FABFAF54222886AE989CFF2D8F0994C812D807529F17D8E468D688398F015533D27CE2CF264F1F6AC8EF368FB3F52C45C50CBB80A63B0CBFED7FEFBBBB3E0DC334124D7736126F6C9AF34D79B9A7FAB77F418966A4C1EBCBA6C9040A1852DF54ED3F8FBB759A575B5BAD539B1F9F8B6C831F6092C83FC384EF316930A12F611273C65FC344FE474C94117C8289D97D808938E86B4CF80B9CD6EF3179F2CC281F70428CCCF6AE9F14234694F811290AEA9C677D832A9AA079C6C09EDEA2D218C84C91D4C91B566414D5717CCB8A25984D17F117B0BC579390DF8BAF6951A08822F97570A1B5FE1558F8112D22FF335A1ADFB5550818292F9CAFD012EFE2B441BF40EBD2870B2D566781F9135A92BF1BED3D2D8C1C5C3EFB42EBB3DE7A2E60D3A76A4CF5B7BA1AF3A7CDF3FF010000FFFFECEF"
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
	println("data length: ", len(double))
	list, err := components.CheckMultiStream(double[:])
	if err != nil {
		t.Error(err)
		return
	}

	println(len(list), err)
	// var dataList [][]byte
	// buf := bytes.NewBuffer([]byte{})
	// read := bytes.NewBuffer(double)
	// finish := false
	// for {
	// 	n, err := read.ReadBytes(0xec)
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Printf("%x ", n)
	// 	fmt.Println(bytes.Equal(n[len(n)-2:], []byte{0xef, 0xec}))
	// 	buf.Write(n)
	// 	if finish && bytes.Equal(n[len(n)-2:], []byte{0xef, 0xec}) {
	// 		fmt.Printf("finish: %x\n", buf.Bytes())
	// 		dataList = append(dataList, buf.Bytes())
	// 		buf = bytes.NewBuffer([]byte{})
	// 		finish = false
	// 		continue
	// 	}

	// 	if len(n) == 2 && !finish {
	// 		finish = true
	// 	}
	// }
	// fmt.Println("data list length:", len(dataList))
	// for i, v := range dataList {
	// 	fmt.Printf("%d:%x\n", i, v)
	// }

	// fmt.Printf("%c\n", 0x34)
	// fmt.Printf("3344:%s", string([]byte{0x33, 0x34}))
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
