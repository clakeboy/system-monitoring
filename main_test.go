package main

import (
	"fmt"
	"github.com/asdine/storm/q"
	"io"
	"system-monitoring/models"
	"testing"
	"time"
)

func TestEmbedFS(t *testing.T) {
	f, err := htmlFiles.Open("assets/html/test.html")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(io.ReadAll(f))
	f.Close()
}

func TestNewNodeInfoModel(t *testing.T) {

	model := models.NewNodeInfoModel("127.0.0.1")
	//var list []models.NodeInfoData
	rangeDay := time.Now().Unix()
	query := model.Select(q.Lt("CreatedDate", rangeDay-(24*3600)))
	num, err := query.Count(new(models.NodeInfoData))
	//num, err := model.Count(new(models.NodeInfoData))
	fmt.Println(num, err)
}
