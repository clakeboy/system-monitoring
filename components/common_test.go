package components

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"math"
	"os"
	"testing"
	"time"
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

func TestGzip(t *testing.T) {
	str := `map[disk0:{"readCount":55099132,"mergedReadCount":0,"writeCount":13479831,"mergedWriteCount":0,"readBytes":1754326626304,"writeBytes":333648203776,"readTime":52569765,"writeTime":10499536,"iopsInProgress":0,"ioTime":63069302,"weightedIO":0,"name":"disk0","serialNumber":"","label":""}]`
	org := []byte(str)
	fmt.Println("original data length:", len(org))
	data, err := Gzip(org)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(" data length:", len(data))
	fmt.Println(string(data))

	unData, err := UnGzip(data)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(unData))
}

func TestBin(t *testing.T) {
	//fmt.Println(utils.ByteToBinaryString(-5))
	bb := []byte("acb我的你是在c是dw是的😀")

	gb, _ := utils.UTF82GBK(bb)
	fmt.Printf("%X\n", bb)
	fmt.Printf("%X\n", gb)
	fmt.Println("\U0001F436")
	fmt.Println(math.Pow(65535, 8))
	fmt.Printf("%08b\n", 69)
	//34024083076439104000000000000000000000
	fmt.Println(string(rune(69)))
}

func TestData(t *testing.T) {
	dbconf := &ckdb.DBConfig{
		DBServer:   "168.168.0.10",
		DBName:     "pcbx_integral_shop_20210716",
		DBUser:     "root",
		DBPassword: "kKie93jgUrn!k",
		DBPort:     "3306",
		DBPoolSize: 200,
		DBIdleSize: 100,
		DBDebug:    true,
	}

	db, _ := ckdb.NewDBA(dbconf)
	sql := `select * from (select usr_name,usr_phone,count(usr_phone) as nums,from_unixtime(create_time) as date_str from t_exchange_record
where create_time > ? and create_time < ?
group by usr_phone) as record where nums > 2`
	date, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-06-01 00:00:00", time.Local)
	var list [][]string
	for i := 1; i <= 48; i++ {
		start := date.Unix()
		res, err := db.Query(sql, start, start+24*3600)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(date.Format("2006-01-02 15:04:05"), ":")
		utils.PrintAny(res)
		date = date.AddDate(0, 0, 1)
		if len(res) > 0 {
			for _, v := range res {
				ress, err := db.Query("select *,from_unixtime(create_time) as rdate from t_exchange_record where usr_phone=?", v["usr_phone"])
				if err != nil {
					fmt.Println("query error", v["usr_phone"])
					continue
				}
				for _, u := range ress {
					var syNo, payFee string
					setInfo, err := db.QueryOne("select * from t_settlement where club_code=?", u["club_code"].(string))
					if err != nil {
						syNo = ""
						payFee = ""
					} else {
						syNo = setInfo["sy_no"].(string)
						payFee = fmt.Sprintf("%d", setInfo["pay_fee"].(int64)/100)
					}

					rs := []string{
						u["usr_name"].(string),
						u["usr_phone"].(string),
						u["club_code"].(string),
						u["card_key"].(string),
						syNo,
						payFee,
						u["rdate"].(string),
					}
					list = append(list, rs)
				}
			}
		}
	}

	fmt.Println(list)
	f, _ := os.Create("./list.csv")
	csw := csv.NewWriter(f)
	csw.Write([]string{
		"用户名",
		"手机",
		"兑换号码",
		"卡号码",
		"保单号",
		"支付金额",
		"兑换时间",
	})
	err := csw.WriteAll(list)
	if err != nil {
		return
	}
	_ = f.Close()
}
