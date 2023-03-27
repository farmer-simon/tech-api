package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/farmer-simon/go-utils"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	AppKey         = "547cbd1f6dfe42caaec325250836d6fb"
	AppSecret      = "d34f1a73f9a14dcd84ecf7c7c3c40a24"
	ApiSuccessCode = "40001"
)

type Client struct {
}

func (c *Client) GetToken() (token string, err error) {
	token, refreshToken, expires := model.CreateSettingsFactory("").GetTokenCache()
	currentTime := time.Now().Unix()
	if token == "" || expires-currentTime < 3600 {
		apiUrl := "http://sc.cqepc.cn:10001/uapservice/open/token/apply"
		postData := url.Values{}
		postData.Set("grantType", "CLIENT_CREDENTIALS")
		client := &http.Client{}
		r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(postData.Encode())) // URL-encoded payload
		r.Header.Add("appKey", AppKey)
		r.Header.Add("appSecret", AppSecret)
		resp, err := client.Do(r)
		if err != nil {
			variable.ZapLog.Error("获取Token失败", zap.Error(err))
			return "", err
		}
		defer resp.Body.Close()
		var apiData map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &apiData); err != nil {
			variable.ZapLog.Error("解析Token数据失败", zap.Error(err))
			return "", err
		}
		if apiData["code"] != ApiSuccessCode {
			variable.ZapLog.Error("Token返回状态错误", zap.Error(errors.New(apiData["message"].(string))))
			return "", errors.New(apiData["message"].(string))
		}
		content := apiData["content"].(map[string]interface{})
		token = content["token"].(string)
		refreshToken = content["refreshToken"].(string)
		expires = currentTime + utils.String2Int64(content["expiresIn"].(string))
		model.CreateSettingsFactory("").SetTokenCache(token, refreshToken, expires)
		return token, nil
	}
	//else if expires-currentTime < 3600 {
	//	token, err = c.refreshToken(refreshToken)
	//	if err != nil {
	//		variable.ZapLog.Error("刷新Token数据失败", zap.Error(err))
	//		return "", err
	//	}
	//}
	return token, nil
}

func (c *Client) refreshToken(refreshToken string) (token string, err error) {
	apiUrl := "http://sc.cqepc.cn:10001/uapservice/open/token/refresh"
	postData := url.Values{}
	postData.Set("grantType", "REFRESH_TOKEN")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(postData.Encode())) // URL-encoded payload
	r.Header.Add("appKey", AppKey)
	r.Header.Add("refreshToken", refreshToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var apiData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &apiData); err != nil {
		return "", err
	}
	fmt.Println(apiData)
	if apiData["code"] != ApiSuccessCode {
		return "", errors.New(apiData["message"].(string))
	}
	currentTime := time.Now().Unix()
	content := apiData["content"].(map[string]interface{})
	token = content["token"].(string)
	refreshToken = content["refreshToken"].(string)
	expires := currentTime + utils.String2Int64(content["expiresIn"].(string))
	model.CreateSettingsFactory("").SetTokenCache(token, refreshToken, expires)
	return token, nil
}

func (c *Client) CheckTicket(ticket string) (userId string, err error) {
	//调试Debug
	if ticket == "cs10000" {
		return "cs10000", nil
	}
	token, err := c.GetToken()
	if err != nil {
		return "", errors.New("获取Token失败，请联系管理员")
	}
	apiUrl := "http://sc.cqepc.cn:10001/uapservice/application/roam/check"
	postData := url.Values{}
	postData.Set("ticket", ticket)
	client := &http.Client{}
	fmt.Println(strings.NewReader(postData.Encode()))
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(postData.Encode())) // URL-encoded payload
	//r.Header.Set("Content-Type", "multipart/form-data")
	r.Header.Add("token", token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var apiData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	//body = []byte("{\"code\":\"40001\",\"message\":\"操作成功\",\"otherMsg\":null,\"content\":{\"userInfo\":{\"userId\":\"cs10000\",\"callApplicationCode\":\"officeHallApplicationCode\"}}}")
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &apiData); err != nil {
		return "", err
	}
	if apiData["code"] != ApiSuccessCode {
		return "", errors.New(apiData["message"].(string))
	}
	fmt.Println(apiData)
	userInfo := apiData["content"].(map[string]interface{})["userInfo"]
	userId = userInfo.(map[string]interface{})["userId"].(string)
	if userId == "" {
		return "", errors.New("授权登录失败，请重试")
	}
	return
}

func (c *Client) SendMessage(userIds []string, message string) error {
	token, err := c.GetToken()
	if err != nil {
		return errors.New("获取Token失败，请联系管理员")
	}

	apiUrl := "http://sc.cqepc.cn:10001/message/appPush/newToDeal"

	type RequestBody struct {
		ApplicationCode string   `json:"applicationCode"`
		UserIds         []string `json:"userIds"`
		Title           string   `json:"title"`
		Classification  string   `json:"classification"`
		PcUrl           string   `json:"pcUrl"`
	}
	//测试期仅发给测试号
	userIds = []string{"cs10000"}

	var rBody = RequestBody{
		ApplicationCode: "h1o5279",
		UserIds:         userIds,
		Title:           message,
		Classification:  "",
		PcUrl:           variable.ConfigYml.GetString("HomePage"),
	}
	if strings.Contains(rBody.PcUrl, "fwei.net") {
		return nil
	}
	sBody, _ := json.Marshal(rBody)

	client := &http.Client{}

	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(string(sBody))) // URL-encoded payload
	//r.Header.Set("Content-Type", "multipart/form-data")
	r.Header.Add("token", token)
	r.Header.Add("Content-Type", "application/json; charset=utf-8")
	r.Header.Add("Content-Length", strconv.Itoa(len(sBody)))
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var apiData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &apiData); err != nil {
		return err
	}
	if apiData["code"] != "1" {
		return errors.New(apiData["error"].(string))
	}

	return nil
}
