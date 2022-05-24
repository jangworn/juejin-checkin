package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	baseUrl          = "https://api.juejin.cn"
	getTodayStatus   = "/growth_api/v1/get_today_status"
	checkInApi       = "/growth_api/v1/check_in"
	getLotteryConfig = "/growth_api/v1/lottery_config/get"
	drawLottery      = "/growth_api/v1/lottery/draw"
)

var cookie string

type Err struct {
	Err_no  int
	Err_msg string
}

type Resp struct {
	Err
	Data bool
}

type Lottery struct {
	Lottery    interface{}
	Free_count int
	point_cost int
}

type RespLottery struct {
	Err
	Data Lottery
}

type RespDraw struct {
	Err
	Data DrawResult
}

type DrawResult struct {
	Lottery_name      string
	Total_lucky_value int
	Draw_lucky_value  int
}

func sendRequest(method string, url string) (result string, err error) {
	var req *http.Request
	if method == "get" {
		req, _ = http.NewRequest("GET", url, nil)
	} else {
		req, _ = http.NewRequest("POST", url, nil)
	}
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(body), nil
}

//签到
func checkIn() {
	ok := getTodayCheckStatus()
	if ok == true {
		fmt.Println("今日已签到,请勿重复签到")
		return
	}
	result, err := sendRequest("post", baseUrl+checkInApi)
	if err != nil {
		fmt.Println("签到失败:", err)
	}
	var r Resp
	json.Unmarshal([]byte(result), &r)
	if r.Err_no > 0 {
		fmt.Println("签到失败:", r.Err_msg)
	} else {
		fmt.Print("签到成功")
	}
}

//查询今日是否已经签到
func getTodayCheckStatus() bool {
	result, err := sendRequest("get", baseUrl+getTodayStatus)
	if err != nil {
		fmt.Println("查询签到失败:", err)
	}
	var r Resp
	json.Unmarshal([]byte(result), &r)
	if r.Err_no == 0 && r.Data == true {
		return true
	} else {
		return false
	}
}

//抽奖
func draw() {
	ok := getTodayDrawStatus()
	if ok == false {
		fmt.Println("免费抽奖次数已用完")
		return
	}
	result, err := sendRequest("post", baseUrl+drawLottery)
	if err != nil {
		fmt.Println("免费抽奖失败:", err)
	}
	var r RespDraw
	json.Unmarshal([]byte(result), &r)
	if r.Err_no > 0 {
		fmt.Println("免费抽奖失败:", r.Err_msg)
	} else {
		fmt.Printf("免费抽奖成功,奖品:%v,增加幸运值:%v,总幸运值:%v", r.Data.Lottery_name, r.Data.Draw_lucky_value, r.Data.Total_lucky_value)
	}
}

//获取今日免费抽奖次数
func getTodayDrawStatus() bool {
	result, err := sendRequest("get", baseUrl+getLotteryConfig)
	if err != nil {
		fmt.Println("查询免费抽奖次数失败:", err)
	}
	var r RespLottery
	json.Unmarshal([]byte(result), &r)
	if r.Err_no == 0 && r.Data.Free_count > 0 {
		return true
	} else {
		return false
	}
}

func main() {
	cookie = os.Args[1]
	if cookie == "" {
		fmt.Println("请传入cookie")
		return
	}
	checkIn()
	draw()
}
