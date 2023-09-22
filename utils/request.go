package utils

import (
	"encoding/json"
	"fmt"
	"github.com/canteen_management/logger"
	"io/ioutil"
	"net/http"
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
