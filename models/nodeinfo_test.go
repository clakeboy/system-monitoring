package models

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/utils"
)

func TestNewNodeInfoModel(t *testing.T) {
	model := NewNodeInfoModel("127.0.0.1")
	//fmt.Println(model.GetById(1))
	count, err := model.Count(new(NodeInfoData))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(count)
	rangeDay := time.Now().Unix()
	// tx, err := model.Begin(true)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// fmt.Println(rangeDay-(24*3600), rangeDay)
	// query := tx.Select(q.Gte("CreatedDate", rangeDay-(24*3600)))
	// num, err := query.Count(new(NodeInfoData))
	// fmt.Println(num)

	list, err := model.List(1, 100, q.Lte("CreatedDate", rangeDay-(2450)))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(list))
	for _, v := range list {
		fmt.Println(time.Unix(v.CreatedDate, 0).Format("2006-01-02T15:04:05Z"))
	}
}

func TestSplit(t *testing.T) {
	reg, _ := regexp.Compile("(\\d{1,2})-\\d{1,2}-\\d{1,2}")
	str := `早             
中             
晚 2-23-3  2-25-4  1-30-3  1-24-3 4-23-1        
早             
中             
晚 13-21-2  6-3-3  14-32-3  14-26-1  1-27-1  9-2-2  12-21-2      
早             
中             
晚 1-23-3  14-26-3  1-19-3  1-29-1  1-18-1        
早 2-27-2  1-19-3  2-19-2  14-7-3         
中 5-7-1  1-6-4  12-9-1  3-14-1  16-4-1  1-19-1        
晚 1-20-2  1-10-4  1-27-2  4-18-2  2-14-1         
早 1-2-1  1-28-3  13-26-3  12-27-4  1-28-2  13-24-3        
中 4-2-1  10-7-1  14-26-3  1-22-2         
晚  2-3-4  1-32-4  2-23-3  4-25-3  4-23-1  4-7-1  14-21-3      
早 6-7-1  1-4-3  12-15-4  1-5-1  14-21-3  10-16-2        
中  1-5-1  10-8-4  14-21-3  12-15-4           
晚 1-27-1  1-23-3  14-21-3  1-29-2  10-10-3  6-3-2        
早 1-5-1  3-14-2  5-2-1  13-12-3  13-26-3  12-20-2        
中 1-4-3  9-3-1  13-27-1          
晚 2-12-3  1-30-3  9-2-2  1-9-2  1-24-3        
早  3-14-2 14-3-1 13-26-3 14-7-3 13-24-3 5-7-1 1-4-3 2-26-3 13-6-2 12-20-2 2-25-4 
中 14-26-3  14-21-3            
晚 10-10-4  4-19-3  11-27-1 1-32-4  14-14-4  14-32-4        
早 14-8-1  12-20-2  1-7-2           
中 12-19-4 1-29-1  1-31-2  2-12-3          
晚  1-20-3 13-21-2 4-17-2 2-3-2 14-32-3 14-26-1       
早 13-27-1  2-19-3  1-13-1           
中 11-6-2  3-8-1 1-4-1 3-4-2         
晚 11-6-2 14-26-3 13-2-4 1-9-1 3-18-2 3-15-1       
早 2-16-1 8-2-1-2 1-12-2 14-12-4 11-16-1        
中 14-4-1 11-6-2 1-31-2 14-12-4 1-12-2        
晚 1-29-4 1-29-1 7-2-2-1 12-10-1 13-4-2        
早 13-15-3  12-18-4 13-27-4 8-1-4-2          
中 1-18-2 2-12-3           
晚 2-28-2 14-12-4 9-2-2 2-27-1         
早 1-21-3 1-2-1 1-4-3 1-28-2 4-13-3        
中 1-7-2 1-7-2           
晚 3-11-1 14-8-1 1-31-4 1-7-2 1-7-2        
早 2-19-3 2-5-2 2-3-2 2-3-3 1-6-4 4-21-2 14-12-4 2-19-2 6-6-1 14-7-3 11-6-2 14-3-1 
中 2-28-3 2-25-4 2-27-2 2-7-4 2-3-4        
晚 2-20-4 2-9-3 2-6-2 2-28-2         
早 1-13-3 4-17-2 4-2-1 4-1-2 4-26-3 4-14-3       
中 4-21-2 4-8-3 4-7-4 4-9-1 4-7-3 4-3-1       
晚 4-3-1 1-16-2 4-16-1 4-22-1         
早 1-5-1 6-18-3 13-4-2 1-19-2         
中 1-19-2 6-18-3 3-4-2          
晚 13-21-2 1-9-1 1-19-3 1-29-1         
早 1-13-1 14-11-1 1-30-3 14-17-1 13-26-3        
中 14-32-3 14-32-1           
晚 1-27-1 1-7-3 1-10-2 14-32-3         
早 12-4-1  12-6-1 12-6-3 12-7-3 12-8-2 2-19-3 12-20-2      
中 12-4-1  12-9-1 12-14-1 12-19-2  12-21-2         
晚 12-4-1  12-9-1  12-10-1 12-21-2 12-24-4        
早 2-19-3            
中 2-5-3 9-3-1           
晚 14-14-2 14-14-4 14-12-4          
早 2-5-3 13-26-3 13-6-2 13-26-4 13-24-3 13-16-2 13-14-1 2-14-3 2-27-2 14-12-4 1-17-3 11-6-2 14-7-3
中 2-5-3 13-15-3 13-27-1 13-12-3 13-7-2        
晚 2-12-3 13-9-1 13-21-2 13-4-2 13-5-1        
早 11-19-1 11-27-4 11-27-1 11-27-1         
中 11-6-2 11-22-4 11-27-2 6-18-3 2-5-3        
晚 11-10-2 11-19-1 11-27-1 11-26-1         
早 3-14-2 1-28-2 6-9-2          
中 2-5-3 1-12-4           
晚 1-29-2 14-14-2 5-3-1 6-9-2         
早 2-5-3 1-17-3 14-12-4          
中 2-5-3 1-17-3 14-12-4          
晚 10-10-4 1-7-2 1-7-2 1-27-1         
早 1-17-3 14-14-2 14-14-2 11-6-2         
中 1-12-2            
晚 14-14-2 1-10-4 9-2-2 14-8-1         
早 6-9-2 1-24-3           
中 11-6-2 11-19-1           
晚 7-2-2-1 14-14-2 4-8-3 3-24-1         
早 14-5-3 1-29-4 2-26-2 4-12-1 4-12-1 13-24-3 2-5-3 1-2？？     
中 1-12-2 2-5-3           
晚 1-9-2 1-9-2 14-28-2 14-28-2         
早 1-9-2 3-4-1 3-14-2 6-9-2         
中 6-9-2            
晚 14-14-2 1-16-2 6-9-2 1-30-3         
早 13-15-3 2-5-3           
中 13-12-3 2-5-3           
晚 14-14-2 3-11-1 13-21-2 1-27-1         
早 12-18-4 6-18-3 1-7-1 4-13-3 2-5-3 6-18-3       
中 2-26-2 1-29-2 2-5-3 6-18-3         
晚 3-18-2 3-18-2 11-3-2 1-24-3         
早 1-9-4 4-13-3           
中             
晚 2-19-3 1-29-2 14-14-2 10-7-3         
早             
中             
晚 1-7-2 1-7-2 9-2-2 1-9-2         
早 11-6-2            
中 11-6-2            
晚             `
	rd := strings.NewReader(str)
	buf := bufio.NewReader(rd)
	count := 0
	treeData := map[string]map[string]map[string]int{}
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}
		subMatch := reg.FindAllSubmatch(line, -1)
		if subMatch == nil {
			continue
		}
		sub := []rune(string(line))
		//fmt.Printf("%s ", string(sub[0]))
		head := string(sub[0])
		count += len(subMatch)
		for _, v := range subMatch {
			buildNo := string(v[1])
			roomNo := string(v[0])
			build, ok := treeData[buildNo]
			if !ok {
				build = make(map[string]map[string]int)
			}

			switch head {
			case "早", "中":
				if _, ok := build["早"]; !ok {
					build["早"] = make(map[string]int)
				}
				build["早"][roomNo] += 1
			case "晚":
				if _, ok := build["晚"]; !ok {
					build["晚"] = make(map[string]int)
				}
				build["晚"][roomNo] += 1
			}

			treeData[buildNo] = build

			fmt.Print(string(v[0]), " ", string(v[1]), " ")
		}
		fmt.Println("")
	}

	fmt.Println("count:", count)

	var buildNos []string
	for k, _ := range treeData {
		buildNos = append(buildNos, k)
	}
	sort.Slice(buildNos, func(i, j int) bool {
		s, _ := strconv.Atoi(buildNos[i])
		e, _ := strconv.Atoi(buildNos[j])
		return s < e
	})
	mo := 0
	ne := 0
	for _, buildNo := range buildNos {
		if bdn, _ := strconv.Atoi(buildNo); bdn > 14 {
			continue
		}
		v := treeData[buildNo]
		fmt.Println("楼栋号：", buildNo)
		bmo := 0
		bne := 0
		for m, d := range v {
			fmt.Printf("    %s班\n", m)
			cc := 0
			for r, co := range d {
				fmt.Printf("        %-9s 次数:%2d %4d\n", r, co, utils.YN(m == "晚", co*150, 100))
				//fmt.Println("     ", r, "次数", co, utils.YN(m == "晚", co*150, ""))
				if m == "晚" {
					cc += co
				}
			}
			if m == "晚" {
				bne = cc * 150
				ne += bne
				fmt.Printf("        晚班总数: %d, 金额: %d 元\n", cc, bne)
			} else {
				bmo = len(d) * 100
				mo += bmo
				fmt.Printf("        白班每户合并为1次 总数: %d, 金额: %d 元\n", len(d), bmo)
			}
		}
		fmt.Printf("    %s栋总金额: %4d\n", buildNo, bmo+bne)
	}

	fmt.Println("\n白班-总发放金额:", mo, "次数:", mo/100)
	fmt.Println("晚班-总发放金额:", ne, "次数:", ne/150)
	fmt.Println("总发放数: ", mo+ne)
}

type BuildNo struct {
	Name       string
	Rooms      []*Room
	Total      float64
	Percentage float64
	Remain     float64
}

type Room struct {
	Label      string
	Level      int
	Number     int
	Fee        float64
	Percentage float64
	Remain     float64
}

func TestStatMoney(t *testing.T) {
	con, err := os.ReadFile("/Users/clakeboy/Downloads/维权募捐明细.txt")
	if err != nil {
		t.Error(err)
		return
	}

	reg, _ := regexp.Compile(`((\d+)-\d+-\d(-\d)?)\s+(\d+)`)
	list := reg.FindAllStringSubmatch(string(con), -1)
	buildList := make(map[string]*BuildNo)
	totalMoney := 0.0
	totalRemain := 31757.0
	for _, v := range list {
		room := &Room{}
		room.Label = v[1]
		fee, err := strconv.ParseFloat(v[4], 64)
		if err != nil {
			fee = 0
		}
		room.Fee = fee

		tmpList := strings.Split(room.Label, "-")
		bno := tmpList[0]
		if len(tmpList) > 3 {
			room.Level, _ = strconv.Atoi(tmpList[2])
			room.Number, _ = strconv.Atoi(tmpList[3])
		} else {
			room.Level, _ = strconv.Atoi(tmpList[1])
			room.Number, _ = strconv.Atoi(tmpList[2])
		}

		build, ok := buildList[bno]
		if !ok {
			build = &BuildNo{
				Name: bno + "栋",
			}
			buildList[bno] = build
		}
		build.Total += fee
		build.Rooms = append(build.Rooms, room)
		totalMoney += fee
	}
	var keys []string
	for k, _ := range buildList {
		keys = append(keys, k)
	}

	//sort.Strings(keys)
	sort.Slice(keys, func(i, j int) bool {
		one, _ := strconv.Atoi(keys[i])
		two, _ := strconv.Atoi(keys[j])
		return one < two
	})

	buildRemain := 0.0
	for _, k := range keys {
		v := buildList[k]
		fmt.Println(v.Name)
		roomsRemain := 0.0
		v.Percentage = v.Total / totalMoney
		v.Remain = math.Round(v.Percentage * totalRemain)
		sort.Slice(v.Rooms, func(i, j int) bool {
			one := v.Rooms[i]
			two := v.Rooms[j]
			if one.Level < two.Level {
				return true
			} else if one.Level == two.Level && one.Number < two.Number {
				return true
			} else {
				return false
			}
		})
		for _, room := range v.Rooms {
			room.Percentage = room.Fee / v.Total
			room.Remain = math.Floor(room.Percentage * v.Remain)
			//room.Percentage = room.Fee / totalMoney
			//room.Remain = math.Floor(room.Percentage * totalRemain)
			roomsRemain += room.Remain
			fmt.Printf("%-8s，捐款：%.f，退款：%.f，退款占楼栋款比：%% %.2f\n", room.Label, room.Fee, room.Remain, room.Percentage*100)
		}
		buildRemain += roomsRemain
		fmt.Printf("楼栋捐款总额：%.2f，每户退款相加：%.2f -- 楼栋比例占总退款相乘：%.2f\n\n", v.Total, roomsRemain, v.Remain)
	}
	fmt.Printf("计算退款总数：%.2f -- 原始退款总数：%.2f\n", buildRemain, totalRemain)
	fmt.Printf("捐款总数: %.2f\n", totalMoney)
	//utils.PrintAny(buildList)
}
