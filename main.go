package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseUrl          = "https://api.juejin.cn"
	getTodayStatus   = "/growth_api/v1/get_today_status"
	checkInApi       = "/growth_api/v1/check_in"
	getLotteryConfig = "/growth_api/v1/lottery_config/get"
	drawLottery      = "/growth_api/v1/lottery/draw"
	cookie           = "MONITOR_WEB_ID=a176c5b8-9a88-468b-aa21-44426d11b9e9; _ga=GA1.2.1104176138.1634093821; _tea_utm_cache_2608={%22utm_source%22:%22timeline_5%22%2C%22utm_medium%22:%22banner%22%2C%22utm_campaign%22:%22xiaoce_linda_20211117%22}; passport_csrf_token_default=8eea6bc27dfccdffe8f8ae97ee992460; passport_csrf_token=8eea6bc27dfccdffe8f8ae97ee992460; sid_guard=9c3dc2094e8c091b9e506e74dc242f55%7C1637634333%7C5184000%7CSat%2C+22-Jan-2022+02%3A25%3A33+GMT; uid_tt=9e6c94dbf429fecba42b7d4250be617e; uid_tt_ss=9e6c94dbf429fecba42b7d4250be617e; sid_tt=9c3dc2094e8c091b9e506e74dc242f55; sessionid=9c3dc2094e8c091b9e506e74dc242f55; sessionid_ss=9c3dc2094e8c091b9e506e74dc242f55; sid_ucp_v1=1.0.0-KDNjNjBlZjVmMDVjMmIwNzViNTc1OWYxYTIyZTEyNjlmZTViZDk3ODgKFwjej5D-hPWXAxCdovGMBhiwFDgCQPEHGgJsZiIgOWMzZGMyMDk0ZThjMDkxYjllNTA2ZTc0ZGMyNDJmNTU; ssid_ucp_v1=1.0.0-KDNjNjBlZjVmMDVjMmIwNzViNTc1OWYxYTIyZTEyNjlmZTViZDk3ODgKFwjej5D-hPWXAxCdovGMBhiwFDgCQPEHGgJsZiIgOWMzZGMyMDk0ZThjMDkxYjllNTA2ZTc0ZGMyNDJmNTU; n_mh=u9ffqAcwOPX6awLXCM_SHqJLoECKsQDK3up-ESBCEdE; _gid=GA1.2.1998507205.1640567014"
)

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
	checkIn()
	draw()
}
