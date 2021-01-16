package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mdl/wlog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserAdd(c *gin.Context) {
	param := &struct {
		Phone   string `json:"phone"`
		SmsCode string `json:"smsCode"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		return
	}

	user := NewUser(param.Phone)

	post := url.Values{}
	post.Set("param", fmt.Sprintf(`{"authCode":"%s","authorized":false,"cid":"%s","phone":"%s"}`, param.SmsCode, user.CID, param.Phone))

	req, err := http.NewRequest(http.MethodPost, URLLogin, strings.NewReader(post.Encode()))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "创建请求失败"})
		return
	}

	req.Header.Add("uuid", user.UUID)
	req.Header.Add("cid", user.CID)
	req.Header.Add("platform", "ANDROID")
	req.Header.Add("apiVersion", "4.9.0")
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("请求登陆时失败:", param.Phone, err)
		c.JSON(http.StatusOK, gin.H{"code": ErrReq, "message": "请求登录时失败"})
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "请求短信验证失败:"})
		return
	}
	loginInfo := &APILogin{}
	if err := json.Unmarshal(data, loginInfo); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrReq, "message": "解析json失败"})
		return
	}

	if loginInfo.Code != "0000" {
		c.JSON(http.StatusOK, gin.H{"code": ErrReq, "message": fmt.Sprintf("登录失败:%s", loginInfo.Result)})
		return
	}

	user.Phone = param.Phone
	user.RealName = loginInfo.Data.WebUser.Realname
	user.UserID = loginInfo.Data.WebUser.ID
	user.TicketName = loginInfo.Data.WebUser.TicketName
	user.Token = loginInfo.Data.WebUser.Token

	if qinfo, err := qualification(user); err != nil {
		wlog.Error("获取用户抢购资格数据失败:", err)
	} else {
		user.Qualification = qinfo.Data.HasQualification
		user.RemainCount = qinfo.Data.RemainCount
	}

	DB.Create(user)

	c.JSON(http.StatusOK, gin.H{"code": NoErr})
}

func SendSms(c *gin.Context) {
	param := &struct {
		Phone string `json:"phone"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "参数解析错误"})
		return
	}

	post := url.Values{}
	post.Set("param", fmt.Sprintf(`{"graphCode":"","isVoiceValidateCode":false,"phone":"%s","type":"login"}`, param.Phone))

	req, err := http.NewRequest(http.MethodPost, URLValidCode, strings.NewReader(post.Encode()))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "创建请求失败"})
		return
	}

	req.Header.Add("uuid", GenUUID())
	req.Header.Add("apiVersion", "4.9.0")
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("请求发送验证码失败:", param.Phone, err)
		c.JSON(http.StatusOK, gin.H{"code": ErrReq, "message": "请求登录时失败"})
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "请求短信验证失败:"})
		return
	}
	validCode := &APIValidCode{}
	if err := json.Unmarshal(data, validCode); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrReq, "message": "解析json失败"})
		return
	}

	if validCode.Code != "0000" {
		c.JSON(http.StatusOK, gin.H{"code": ErrReq, "message": fmt.Sprintf("发送短信失败:%s", validCode.Msg)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": NoErr})
}

func SetStore(c *gin.Context) {
	param := &struct {
		Phone    string `json:"phone"`
		StoreIDs []int  `json:"storeIDs"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "参数解析错误"})
		return
	}

	if param.Phone == "" {
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "手机号不能为空"})
		return
	}

	strStoreIDs := ""
	for _, storeID := range param.StoreIDs {
		if len(strStoreIDs) == 0 {
			strStoreIDs = strconv.Itoa(storeID)
		} else {
			strStoreIDs = fmt.Sprintf("%s;%d", strStoreIDs, storeID)
		}
	}

	DB.Model(&TableUser{}).Where(&TableUser{Phone: param.Phone}).Update("store_id", strStoreIDs)

	c.JSON(http.StatusOK, gin.H{"code": NoErr})
}

func UserList(c *gin.Context) {
	param := &struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}{Page: 1, Limit: 20}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "参数解析错误"})
		return
	}

	totalNum := int64(0)
	users := make([]*struct {
		*TableUser
		StoreName string `json:"store_name" gorm:"-"`
	}, 0)

	result := DB.Model(&TableUser{}).
		Count(&totalNum).
		Offset((param.Page - 1) * param.Limit).
		Limit(param.Limit).
		Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrDB, "message": "数据库读取错误"})
		return
	}

	for _, user := range users {
		if user.StoreID == "" {
			continue
		}

		stores := strings.Split(user.StoreID, ";")
		for _, singleStore := range stores {
			if singleStore == "" {
				continue
			}
			id, err := strconv.Atoi(singleStore)
			if err != nil {
				wlog.Error("商店id转换失败:", err)
				continue
			}
			if storeInfo, ok := StoreMap[id]; ok {
				user.StoreName = fmt.Sprintf("%s;%s", user.StoreName, storeInfo.StoreName)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": NoErr, "users": users, "totalNum": totalNum})
}

func UserDelete(c *gin.Context) {
	param := &struct {
		Phone string `json:"phone"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "参数解析错误"})
		return
	}

	DB.Where("phone = ?", param.Phone).Delete(&TableUser{}).Limit(1)

	c.JSON(http.StatusOK, gin.H{"code": NoErr})
}

func UserEdit(c *gin.Context) {

}

func OrderList(c *gin.Context) {
	param := &struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}{Page: 1, Limit: 20}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		c.JSON(http.StatusOK, gin.H{"code": ErrParam, "message": "参数解析错误"})
		return
	}

	totalNum := int64(0)
	orders := make([]*TableOrder, 0)

	result := DB.
		Count(&totalNum).
		Offset((param.Page - 1) * param.Limit).
		Limit(param.Limit).
		Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": ErrDB, "message": "数据库读取错误"})
		return
	}
}

func GetStoreList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": NoErr, "list": StoreList})
}
