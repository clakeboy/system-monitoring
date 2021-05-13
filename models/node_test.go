package models

import (
	"bytes"
	"fmt"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"github.com/elastic/go-elasticsearch/v7"
	"testing"
)

func TestESData(t *testing.T) {
	dbconf := &ckdb.DBConfig{
		DBServer:   "168.168.2.12",
		DBName:     "pcbx_nk",
		DBUser:     "root",
		DBPassword: "123123",
		DBPort:     "3306",
		DBPoolSize: 200,
		DBIdleSize: 100,
		DBDebug:    true,
	}

	db, _ := ckdb.NewDBA(dbconf)
	tb := db.Table("t_policy")
	res, err := tb.Where(utils.M{"id": 401}, "").Query().Find()
	logErr(err)
	utils.PrintAny(res)
	fmt.Printf("%+t", res["id"])
}

func getDBA() *ckdb.DBA {
	dbconf := &ckdb.DBConfig{
		DBServer:   "168.168.2.12",
		DBName:     "pcbx_nk",
		DBUser:     "root",
		DBPassword: "123123",
		DBPort:     "3306",
		DBPoolSize: 200,
		DBIdleSize: 100,
		DBDebug:    true,
	}

	db, _ := ckdb.NewDBA(dbconf)
	return db
}

func TestCreateIndex(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://168.168.2.21:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	logErr(err)

	db := getDBA()

	tb := db.Table("t_policy")
	run := utils.NewExecTime()
	run.Start()
	pageNum := 100
	for i := 1; i <= 5; i++ {
		list, err := tb.Limit(pageNum, i).Query().Result()
		fmt.Println(db.LastSql)
		logErr(err)
		utils.PrintAny(list[0])
		createIndex(list, es)
		fmt.Println("add page", i, "number:", i*pageNum)
	}
	run.End(true)
}

func createIndex(list []utils.M, es *elasticsearch.Client) {
	buf := new(bytes.Buffer)
	for _, v := range list {
		//row := v.(utils.M)
		header := utils.M{
			"index": utils.M{
				"_index": "pcbx_nk_1",
				"_id":    v["id"].(string),
			},
		}
		buf.WriteString(header.ToJsonString() + "\n")
		buf.WriteString(v.ToJsonString() + "\n")
	}

	res, err := es.Bulk(buf)
	if err != nil {
		logErr(err)
	}

	if res.IsError() {
		fmt.Println(res.Status())
	}

	res.Body.Close()
}

func logErr(err error) {
	if err != nil {
		panic(err)
		return
	}
}
