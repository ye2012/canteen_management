package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/canteen_management/logger"
)

const (
	wxLoginUrlTemplate = "https://api.weixin.qq.com/sns/jscode2session?appid=%v&secret=%v&js_code=%v&grant_type=authorization_code"

	requestLogTag = "Request"
)

type MiniProgramLoginRes struct {
	Code       int32  `json:"errcode"`
	Message    string `json:"errmsg"`
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
}

func MiniProgramLogin(appID, secret, code string) (string, error) {
	wxLoginUrl := fmt.Sprintf(wxLoginUrlTemplate, appID, secret, code)
	logger.Debug(requestLogTag, "wxLoginUrl:%v", wxLoginUrl)
	resp, err := http.Get(wxLoginUrl)
	if err != nil {
		logger.Warn(requestLogTag, "Request WxLogin Failed|Err:%v", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := &MiniProgramLoginRes{}
	err = json.Unmarshal(body, res)
	if err != nil {
		logger.Warn(requestLogTag, "WxLogin Unmarshal Failed|Err:%v", err)
		return "", err
	}
	if res.Code != 0 {
		logger.Warn(requestLogTag, "WxLogin Response Failed|Err:%v", err)
		return "", fmt.Errorf("login failed|msg:%v", res.Message)
	}
	return res.OpenID, nil
}

func GenerateSign(paramMap map[string]string, signSecret string) string {
	delete(paramMap, "sign")

	key := make([]string, 0)
	for k := range paramMap {
		key = append(key, k)
	}
	sort.Strings(key)

	paramStr := &bytes.Buffer{}
	for _, k := range key {
		paramStr.WriteString(k)
		paramStr.WriteString(fmt.Sprintf("%v", paramMap[k]))
	}
	paramStr.WriteString(signSecret)

	// 使用MD5对待签名串求签
	return GetMD5Hex(paramStr.String())
}
