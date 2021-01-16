package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mdl/wlog"
	"net/http"
	"net/url"
	"strings"
)

func buy(user *TableUser) bool {
	// 进入队列
	for {
		if !jumpQueue(user) {
			continue
		}

		if !QueryResult(user) {
			continue
		}

		// 插入记录
		return true
	}
}

func qualification(user *TableUser) (*APIQualification, error) {
	post := url.Values{}
	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"%s","skuId":"1003057062"}`, user.UUID))

	req, err := http.NewRequest(http.MethodPost, URLQualication, strings.NewReader(post.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", fmt.Sprintf("ticketName=%s", user.TicketName))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	qInfo := &APIQualification{}
	if err := json.Unmarshal(data, qInfo); err != nil {
		return nil, err
	}

	if qInfo.Code != "0000" {
		return nil, errors.New(qInfo.Code)
	}

	return qInfo, nil
}

// 库存查询
func queryStock(storeID int) (*APIStockInfo, error) {
	token := "ba5c2cd7-6daa-4382-b3a6-0213a08db075"
	ticketName := "0CD7A6271D33EEA145A4796F2039915DA3DA30CCCCAE737DBA3787B6A6BC9B95A67273D595C54CF83311CE22354E888BD8880B6BD804718ECA2FA3A47CF04C32EB0A95AAB3EB9A833D9CA5AA8BB19120406402DF04604E65765216560C7C492EF98BCE86C1802B67955327426808FCF727E300860AB89C5CF514DA92DE3FA2B7"
	// 如果用户有多个选择,那么一个个尝试
	// for _, storeID := range storeIDs {
	if storeID == 0 {
		return nil, errors.New("storeID不能为空")
	}

	post := url.Values{}
	post.Set("param", fmt.Sprintf(`{"vendorId":"86","erpStoreId":%d,"skuId":"1003057062","specialMark":1}`, storeID))

	req, err := http.NewRequest(http.MethodPost, URLQueryStock, strings.NewReader(post.Encode()))
	if err != nil {
		return nil, err
	}

	cookie := fmt.Sprintf("venderId=86; businessCode=1; appMode=online; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; token=%s; ticketName=%s; userId=271843147; recommend=1; dmTenantId=2", token, ticketName)
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	stock := &APIStockInfo{}
	if err := json.Unmarshal(data, stock); err != nil {
		return nil, err
	}

	if stock.Code != "0000" {
		return nil, errors.New(stock.Msg)
	}

	return stock, nil
}

func jumpQueue(user *TableUser) bool {
	link := "https://presale.dmall.com/maotai/jumpQueue"
	post := url.Values{}

	// token := "fdd43e4f-0e32-4513-9898-d83746a774ac"
	// ticketName := "4A2F60EDDC26311629E85736F0CD63DCEF3492BFC01ED6C10A1AF315DEEAB3EB1E5D0CDF4A25410B250EB40576A047B01C26FA31032DA689C3CEC392DC5977B0EC7A4CB7FC21EFD41F33113CB3542C917E52ACB5ED3247BE2CA7B79FE115F4E57B47ADAEBDBF2E644D3A55034BDD1C4A706F32475C1602E011E185260BCB3A21"
	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"%s","skuId":"1003057062"}`, user.UUID))
	post.Set("token", user.Token)
	post.Set("ticketName", user.TicketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	cookie := fmt.Sprintf("uuid=%s; venderId=86; businessCode=1; appMode=online; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; token=%s; ticketName=%s; userId=271843147; recommend=1; dmTenantId=2", user.UUID, user.Token, user.TicketName)
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("请求登陆时失败:", err)
		return false
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	jumpQueue := &APIQueueInfo{}
	if err := json.Unmarshal(data, jumpQueue); err != nil {
		wlog.Error("解析失败")
		return false
	}

	if jumpQueue.Code != "0000" {
		wlog.Error("错误:", jumpQueue.Code, jumpQueue.Msg)
		return false
	}

	return true
}

func QueryResult(user *TableUser) bool {
	link := "https://presale.dmall.com/maotai/queryQueue"
	post := url.Values{}
	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"%s","skuId":"1003057062"}`, user.UUID))
	post.Set("token", user.Token)
	post.Set("ticketName", user.TicketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	cookie := fmt.Sprintf("uuid=%s; venderId=86; businessCode=1; appMode=online; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; token=%s; ticketName=%s; userId=271843147; recommend=1; dmTenantId=2", user.UUID, user.Token, user.TicketName)
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("请求登陆时失败:", err)
		return false
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wlog.Error("缓冲读取失败", err)
		return false
	}
	jumpQueue := &APIQueueInfo{}
	if err := json.Unmarshal(data, jumpQueue); err != nil {
		wlog.Error("解析失败", err)
		return false
	}

	if jumpQueue.Code != "0000" {
		wlog.Error("错误:", jumpQueue.Code, jumpQueue.Msg)
		return false
	}

	if jumpQueue.Data {
		wlog.Info("进入队列:", user.Phone, jumpQueue.Msg)
		return true
	}

	return false
}

func TmpPurchase() {
	for {
		if TmpJumpQueue() {
			if TmpQueryResult() {
				wlog.Info("进入队列成功")
				break
			} else {
				wlog.Info("进入队列失败")
			}
		} else {
			wlog.Info("还不能进行排队")
		}
	}
}

func TmpJumpQueue() bool {
	link := "https://presale.dmall.com/maotai/jumpQueue"
	post := url.Values{}

	token := "fdd43e4f-0e32-4513-9898-d83746a774ac"
	ticketName := "4A2F60EDDC26311629E85736F0CD63DCEF3492BFC01ED6C10A1AF315DEEAB3EB1E5D0CDF4A25410B250EB40576A047B01C26FA31032DA689C3CEC392DC5977B0EC7A4CB7FC21EFD41F33113CB3542C917E52ACB5ED3247BE2CA7B79FE115F4E57B47ADAEBDBF2E644D3A55034BDD1C4A706F32475C1602E011E185260BCB3A21"
	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"f549f16478802fb0","skuId":"1003057062"}`))
	post.Set("token", token)
	post.Set("ticketName", ticketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	cookie := fmt.Sprintf("tempid=C935D7076EB00002AAC8C0D11F3116C0; updateTime=1609250551000; inited=true; console_mode=0; web_session_count=1; uuid=f549f16478802fb0; store_id=15535; vender_id=86; storeId=15535; venderId=86; businessCode=1; appMode=online; storeGroupKey=41616635f2585e396aff9eafd786743a@MS0xNTUzNS04Ng; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; token=%s; ticketName=%s; userId=271843147; %s", token, ticketName, `addr=%E6%B9%96%E5%8D%97%E7%9C%81%E9%95%BF%E6%B2%99%E5%B8%82%E5%BC%80%E7%A6%8F%E5%8C%BA%E5%9B%9B%E6%96%B9%E5%9D%AA; community=%E5%9B%9B%E6%96%B9%E5%9D%AA; areaId=430105; session_id=7a13025c68be4feaab3ce195762bc546; env=app; first_session_time=1610424160765; session_count=3; recommend=1; dmTenantId=2`)
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("请求登陆时失败:", err)
		return false
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	jumpQueue := &APIQueueInfo{}
	if err := json.Unmarshal(data, jumpQueue); err != nil {
		wlog.Error("解析失败", string(data), err)
		return false
	}

	if jumpQueue.Code != "0000" {
		wlog.Error("错误:", jumpQueue.Code, jumpQueue.Msg)
		return false
	}

	return jumpQueue.Data
}

func TmpQueryResult() bool {
	link := "https://presale.dmall.com/maotai/queryQueue"
	post := url.Values{}
	token := "fdd43e4f-0e32-4513-9898-d83746a774ac"
	ticketName := "4A2F60EDDC26311629E85736F0CD63DCEF3492BFC01ED6C10A1AF315DEEAB3EB1E5D0CDF4A25410B250EB40576A047B01C26FA31032DA689C3CEC392DC5977B0EC7A4CB7FC21EFD41F33113CB3542C917E52ACB5ED3247BE2CA7B79FE115F4E57B47ADAEBDBF2E644D3A55034BDD1C4A706F32475C1602E011E185260BCB3A21"

	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"f549f16478802fb0","skuId":"1003057062"}`))
	post.Set("token", token)
	post.Set("ticketName", ticketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	cookie := fmt.Sprintf("tempid=C935D7076EB00002AAC8C0D11F3116C0; updateTime=1609250551000; inited=true; console_mode=0; web_session_count=1; uuid=f549f16478802fb0; store_id=15535; vender_id=86; storeId=15535; venderId=86; businessCode=1; appMode=online; storeGroupKey=41616635f2585e396aff9eafd786743a@MS0xNTUzNS04Ng; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; token=%s; ticketName=%s; userId=271843147; %s", token, ticketName, `addr=%E6%B9%96%E5%8D%97%E7%9C%81%E9%95%BF%E6%B2%99%E5%B8%82%E5%BC%80%E7%A6%8F%E5%8C%BA%E5%9B%9B%E6%96%B9%E5%9D%AA; community=%E5%9B%9B%E6%96%B9%E5%9D%AA; areaId=430105; session_id=7a13025c68be4feaab3ce195762bc546; env=app; first_session_time=1610424160765; session_count=3; recommend=1; dmTenantId=2`)
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("请求登陆时失败:", err)
		return false
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wlog.Error("缓冲读取失败", err)
		return false
	}
	jumpQueue := &APIQueueInfo{}
	if err := json.Unmarshal(data, jumpQueue); err != nil {
		wlog.Error("解析失败", string(data), err)
		return false
	}

	if jumpQueue.Code != "0000" {
		wlog.Error("错误:", jumpQueue.Code, jumpQueue.Msg)
		return false
	}

	if jumpQueue.Data {
		wlog.Info("抢购成功:", jumpQueue.Msg)
		return true
	}

	return false
}
