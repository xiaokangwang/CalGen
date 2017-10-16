package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)
import "github.com/davecgh/go-spew/spew"

func main() {
	var sch Schtool
	f, _ := os.Open("in3.csv")
	csvr := csv.NewReader(f)
	r, _ := csvr.ReadAll()
	sch.LoadDem(r)
	sch.SetPtm(true, "4-Sep-2017", 17)
	sch.Analy()
	o := sch.GenerateOut()
	fo, _ := os.Create("out3.csv")
	csvw := csv.NewWriter(fo)
	csvw.WriteAll(o)
	fo.Close()

}

type Schtool struct {
	class       [][]string
	reverseWeek bool
	startTime   string
	classA      [][]class
	maxweek     int
}

func transpose(a [][]string) {
	n := len(a)
	n2 := len(a[0])
	b := make([][]string, n2)
	for i := 0; i < n2; i++ {
		b[i] = make([]string, n)
		for j := 0; j < n; j++ {
			b[i][j] = a[j][i]
		}
	}
	copy(a, b)
}

type class struct {
	name, teacherName string
	stweek            int
	enweek            int
	even, odd, stub   bool
	location          string
	slot              int
}

func (st *Schtool) LoadDem(class [][]string) {
	//spew.Dump(class)
	transpose(class)
	//spew.Dump(class)
	st.class = class
}
func (st *Schtool) SetPtm(reverseWeek bool, startTime string, maxweek int) {
	st.reverseWeek = reverseWeek
	st.startTime = startTime
	st.maxweek = maxweek
}
func (st *Schtool) Analy() {
	resu := make([][]class, 0, 7)
	for _, item := range st.class {
		resur := make([]class, 0, 7)
		for slotn, itemx := range item {
			classs, classs2 := st.genClassFromDescrib(itemx)

			classs.slot = slotn
			classs2.slot = slotn

			resur = append(resur, classs, classs2)
		}
		resu = append(resu, resur)
	}
	st.classA = resu
}
func (st *Schtool) GenerateOut() [][]string {
	headraw := strings.Split("Subject,Start Date,Start Time,End Date,End Time,Private,All Day Event,Location", ",")
	rescand := make([][]string, 0)
	rescand = append(rescand, headraw)

	//weeks loop
	var currentWeek = 1
	var baset, _ = time.ParseInLocation("2-Jan-2006", st.startTime, time.FixedZone("CST", int((time.Hour*8).Seconds())))
	//week
	for currentWeek <= st.maxweek {
		var even = currentWeek%2 == 0
		var odd = currentWeek%2 == 1
		//day
		for index, item := range st.classA {
			var baset2 = baset.Add(time.Hour * time.Duration(24*int64(index)))
			for _, item2 := range item {
				if (even && item2.even) || (odd && item2.odd) {
					nxr := make([]string, 0)
					st1 := fmt.Sprintf("%v|%v", item2.name, item2.teacherName)
					nxr = append(nxr, st1)
					StT, stTE := getTimeWithNum(item2.slot, baset2)
					a, b := genPairFromTime(StT)
					nxr = append(nxr, a, b)
					c, d := genPairFromTime(stTE)
					nxr = append(nxr, c, d)
					nxr = append(nxr, "False", "False")
					nxr = append(nxr, item2.location)
					rescand = append(rescand, nxr)
					//spew.Dump(nxr)
				}
			}
		}
		currentWeek++
		baset = baset.Add(24 * 7 * time.Hour)
	}
	return rescand
}
func genPairFromTime(t time.Time) (string, string) {
	a := t.Format("1/2/2006")
	b := t.Format("3:04 PM")
	return a, b
}
func getTimeWithNum(seq int, t time.Time) (time.Time, time.Time) {
	var thisTS, thisTE time.Time
	switch seq {
	case 0:
		thisTS = t.Add(time.Hour * 8)
		thisTE = t.Add(time.Hour * 9).Add(time.Minute * 50)
	case 1:
		thisTS = t.Add(time.Hour * 10).Add(time.Minute * 10)
		thisTE = t.Add(time.Hour * 12)
	case 2:
		thisTS = t.Add(time.Hour * 13).Add(time.Minute * 30)
		thisTE = t.Add(time.Hour * 15).Add(time.Minute * 20)
	case 3:
		thisTS = t.Add(time.Hour * 15).Add(time.Minute * 40)
		thisTE = t.Add(time.Hour * 17).Add(time.Minute * 30)
	case 4:
		thisTS = t.Add(time.Hour * 18).Add(time.Minute * 00)
		thisTE = t.Add(time.Hour * 21).Add(time.Minute * 00)
	}
	return thisTS, thisTE
}

func (st *Schtool) genClassFromDescrib(describ string) (class, class) {
	if len(describ) <= 5 {
		var cuc class
		cuc.stub = true
		return cuc, cuc
	}
	cur := strings.Split(describ, "\n")
	//spew.Dump(cur)
	cur = cur[1:len(cur)]
	//spew.Dump(cur)
	var cuc class
	cuc.name = cur[0]
	actualname := strings.Split(cur[1], "(")[0]
	cuc.teacherName = actualname
	cuc.location = cur[3]
	//weeks
	var re = regexp.MustCompile(`((?:[0-9])+)-((?:[0-9])+)(?:\[|\()([单双])?周`)
	result := re.FindStringSubmatch(cur[2])
	spew.Dump("describ", describ)
	spew.Dump("Result", cur[2])
	cuc.stweek, _ = strconv.Atoi(result[1])
	cuc.enweek, _ = strconv.Atoi(result[2])
	switch result[3] {
	case "单":
		cuc.odd = true
	case "双":
		cuc.even = true
	default:
		cuc.odd = true
		cuc.even = true
	}

	//spew.Dump(cuc)
	//TODO: Guessed value
	if len(cur) >= 9 {
		cur = cur[4:len(cur)]
		classsss, _ := st.genClassFromDescrib(strings.Join(cur, "\n"))
		return cuc, classsss
	}
	return cuc, class{stub: true}
}
