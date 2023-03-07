package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-resty/resty/v2"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args

	first, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatalln(err)
	}
	last, err := strconv.Atoi(args[4])
	if err != nil {
		log.Fatalln(err)
	}
	step := randomNum(first, last)
	log.Println("Commit step", step)
	loginRequest := LoginRequest{}
	loginRequest.LoginName = args[1]
	loginRequest.Password = fmt.Sprintf("%x", md5.Sum([]byte(args[2])))
	loginRequest.AppType = 6
	loginRequest.ClientID = "88888"
	loginRequest.RoleType = 0
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Dalvik/2.1.0 (Linux; U; Android 5.1.1; Magic2 Build/LMY48Z)",
	}

	loginBody := requestPost("https://sports.lifesense.com/sessions_service/login?systemType=2&version=4.6.7", loginRequest, headers)
	loginObject := LoginResp{}
	loginErr := json.Unmarshal(loginBody, &loginObject)
	log.Println("loginBody", string(loginBody))
	if loginErr != nil {
		log.Fatalln("loginErr", loginErr)
	} else {
		stepRequest := StepRequest{}
		list := StepRequestList{}
		list.DataSource = 2
		list.Active = 1
		list.Calories = strconv.Itoa(int(step / 4)) ////
		list.Distance = int(step / 3)               ////
		list.Step = step                            ///
		list.DataSource0 = 2
		list.DeviceID = "M_NULL"
		list.ExerciseTime = 0
		list.IsUpload = 0
		list.MeasurementTime = time.Now().Format("2006-01-02 15:04:05")
		list.ExerciseTime = 0
		list.Priority = 0
		list.Type = 2
		list.Updated = time.Now().UnixMicro()
		list.UserID, _ = strconv.ParseInt(loginObject.Data.UserID, 10, 64)
		stepRequest.List = append(stepRequest.List, list)
		StepHeader := map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "Dalvik/2.1.0 (Linux; U; Android 5.1.1; Magic2 Build/LMY48Z)",
			"Cookie":       "accessToken=" + loginObject.Data.AccessToken,
		}
		workBody := requestPost("https://sports.lifesense.com/sport_service/sport/sport/uploadMobileStepV2?version=4.5&systemType=2", stepRequest, StepHeader)
		log.Println(string(workBody))
		stepResponse := StepResponse{}
		stepErr := json.Unmarshal(workBody, &stepResponse)
		if stepErr != nil {
			log.Println(stepErr)
		} else {
			log.Println(stepResponse.Msg)
		}
	}

}

func randomNum(first int, last int) (random int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(last-first) + first
	fmt.Println(n)
	return n
}

type StepResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type StepRequest struct {
	List []StepRequestList `json:"list"`
}
type StepRequestList struct {
	DataSource      int    `json:"DataSource"`
	Active          int    `json:"active"`
	Calories        string `json:"calories"`
	DataSource0     int    `json:"dataSource"`
	DeviceID        string `json:"deviceId"`
	Distance        int    `json:"distance"`
	ExerciseTime    int    `json:"exerciseTime"`
	IsUpload        int    `json:"isUpload"`
	MeasurementTime string `json:"measurementTime"`
	Priority        int    `json:"priority"`
	Step            int    `json:"step"`
	Type            int    `json:"type"`
	Updated         int64  `json:"updated"`
	UserID          int64  `json:"userId"`
}

type LoginRequest struct {
	AppType   int    `json:"appType"`
	ClientID  string `json:"clientId"`
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
	RoleType  int    `json:"roleType"`
}

type LoginResp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data LoginRespData `json:"data"`
}
type LoginRespData struct {
	Exist       bool   `json:"exist"`
	HasMobile   bool   `json:"hasMobile"`
	HasEmail    bool   `json:"hasEmail"`
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
	ExpireAt    int64  `json:"expireAt"`
	UserType    int    `json:"userType"`
	NeedInfo    bool   `json:"needInfo"`
}

//发起请求
func requestPost(url string, body interface{}, header map[string]string) (RespBody []byte) {
	request := resty.New().R()
	request.SetHeaders(header)
	request.SetBody(body)
	resp, err := request.Post(url)
	if err != nil {
		return nil
	}
	return resp.Body()
}
