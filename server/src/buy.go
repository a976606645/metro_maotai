package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mdl/wlog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func buy(user *TableUser) {
	wlog.Infof("用户[%s][%s]开始抢购", user.Phone, user.RealName)

	for {
		if !isRun {
			break
		}
		if !jumpQueue(user) {
			continue
		}
		wlog.Infof("用户[%s][%s]开始排队,等待20秒", user.Phone, user.RealName)
		time.Sleep(20 * time.Second)
		if !QueryResult(user) {
			wlog.Infof("用户[%s][%s]没能进入队列,重试", user.Phone, user.RealName)
			continue
		}
		wlog.Infof("用户[%s][%s]进入队列, 进行下单操作", user.Phone, user.RealName)
		if orderSubmit(user) {
			break
		}
	}
}

func orderSubmit(user *TableUser) bool {
	link := "https://presale.dmall.com/maotai/orderSubmit"
	post := url.Values{}

	y, m, d := time.Now().Add(24 * time.Hour).Date()
	shipment := fmt.Sprintf("%d-%d-%d", y, m, d)
	// 商店列表
	stores := strings.Split(user.StoreID, ";")
	for _, singleStore := range stores {
		if singleStore == "" {
			continue
		}
		// storeRemain
		storeID, err := strconv.Atoi(singleStore)
		if err != nil {
			wlog.Error("商店代码转换失败:", singleStore, err)
			continue
		}

		stockInfo, err := queryStock(storeID)
		if err != nil {
			wlog.Error("获取库存失败:", err)
		}
		storeStock := stockInfo.Data.WareInfo.Stock
		if storeStock == 0 {
			wlog.Infof("[%s]购买时[%s]已经没有库存", user.Phone, StoreMap[storeID].StoreName)
			continue
		}

		skuCount := user.RemainCount
		if skuCount > storeStock {
			skuCount = storeStock
		}
		price := skuCount * 149900
		normalParam := `{"vendorId":"86","shipmentStartTime":"09:30","shipmentEndTime":"12:00","skuId":1003057062,"longitude":"","latitude":"","specialMark":1,"orderOrigin":6,"appName":"com.dm.metro","dmTenantId":"2"`
		formatParam := fmt.Sprintf(`%s, "name":"%s", "phone":"%s", "shipmentDate":"%s","erpStoreId":%s, "deviceId":"%s", "skuCount": %d, "price": %d}`, normalParam, user.RealName, user.Phone, shipment, singleStore, user.UUID, skuCount, price)

		post.Set("param", formatParam)
		post.Set("token", user.Token)
		post.Set("ticketName", user.TicketName)

		req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
		if err != nil {
			wlog.Error("创建请求失败:", err)
			return false
		}

		req.Header.Add("cookie", genCookie(user))
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
			wlog.Error("从resp.Body中读取失败", err)
			return false
		}

		result := &APISubmitOrder{}
		if err := json.Unmarshal(data, result); err != nil {
			wlog.Error("解析失败")
			return false
		}

		if result.Code == "1000" {
			wlog.Infof("[%s]抢购成功", user.Phone)

			DB.Create(&TableOrder{
				Phone:   user.Phone,
				Num:     skuCount,
				StoreID: storeID,
			})

			return true

		} else {
			wlog.Infof("[%s]在商店[%s]没有抢到,原因是:%s", user.Phone, StoreMap[storeID].StoreName, result.Msg)
		}
	}

	return false
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
	if storeID == 0 {
		return nil, errors.New("storeID不能为空")
	}

	post := url.Values{}
	post.Set("param", fmt.Sprintf(`{"vendorId":"86","erpStoreId":%d,"skuId":"1003057062","specialMark":1}`, storeID))

	req, err := http.NewRequest(http.MethodPost, URLQueryStock, strings.NewReader(post.Encode()))
	if err != nil {
		return nil, err
	}

	token := "ba5c2cd7-6daa-4382-b3a6-0213a08db075"
	ticketName := "0CD7A6271D33EEA145A4796F2039915DA3DA30CCCCAE737DBA3787B6A6BC9B95A67273D595C54CF83311CE22354E888BD8880B6BD804718ECA2FA3A47CF04C32EB0A95AAB3EB9A833D9CA5AA8BB19120406402DF04604E65765216560C7C492EF98BCE86C1802B67955327426808FCF727E300860AB89C5CF514DA92DE3FA2B7"
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

func genCookie(user *TableUser) string {
	common := "sysVersion=Android-6.0.1; apiVersion=4.9.0; dSource=; oaid=; tdc=; utmId=; channelId=dm010205000001;  platformStoreGroup=; lastInstallTime=1579506068; version=4.9.0; tpc=; firstInstallTime=1579506068; networkType=1; storeGroupV4=; xyz=ac; appName=com.dm.metro; smartLoading=1; utmSource=; wifiState=1; gatewayCache=; platformStoreGroupKey=; isOpenNotification=1; inited=true; console_mode=0; web_session_count=1; vender_id=86; venderId=86; businessCode=1; appMode=online; storeGroupKey=43c9247c2299f05b6551528fbbf23e89@MS04NTA2Mi04Ng; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; userId=206579381; lat=30.706703; lng=111.302199; env=app; first_session_time=1610966352620; session_count=1; recommend=1; dmTenantId=2"

	stores := strings.Split(user.StoreID, ";")
	storeIndex := rand.Intn(len(stores))
	storeID := stores[storeIndex]

	return fmt.Sprintf("tempid=%s; updateTime=%d;androidId=%s; cid=%s; uuid=%s; token=%s; ticketName=%s; device=%s; currentTime=%d; sessionId=%s; User-Agent=%s; store_id=%s; storeId=%s; %s",
		user.TempID,
		user.CreateTime.Unix(),
		user.UUID,
		user.CID,
		user.UUID,
		user.Token,
		user.TicketName,
		user.Device,
		time.Now().Unix(),
		user.SessionID,
		DefaultUserAgent,
		storeID,
		storeID,
		common)
}
func jumpQueue(user *TableUser) bool {
	link := "https://presale.dmall.com/maotai/jumpQueue"
	post := url.Values{}

	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"%s","skuId":"1003057062"}`, user.UUID))
	post.Set("token", user.Token)
	post.Set("ticketName", user.TicketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	req.Header.Add("cookie", genCookie(user))
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

	return jumpQueue.Data
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

	req.Header.Add("cookie", genCookie(user))
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

	return jumpQueue.Data
}

func TmpPurchase() {
	for {
		if TmpJumpQueue() {
			time.Sleep(21 * time.Second)
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

var (
	tempCookie    = "tempid=C937DC44B5400002213A1F1A1BAE105C; updateTime=1610788959000; device=Android%20MuMu%20V417IR%20release-keys; sysVersion=Android-6.0.1; screen=1440*810; apiVersion=4.9.0; dSource=; oaid=; tdc=; utmId=; androidId=8454e52abd5eb0de; channelId=dm010205000001; currentTime=1610966542889; platformStoreGroup=; lastInstallTime=1610424150137; version=4.9.0; tpc=a_146591; firstInstallTime=1610424150137; networkType=1; deliveryLng=; deliveryLat=; cid=52fdfc072182654f163; storeGroupV4=; sessionId=49dd60c7c9f84cebbdddd4ca4b4486f3; User-Agent=dmall/4.9.0%20Dalvik/2.1.0%20%28Linux%3B%20U%3B%20Android%206.0.1%3B%20MuMu%20Build/V417IR%29; xyz=ac; appName=com.dm.metro; smartLoading=1; utmSource=; wifiState=1; gatewayCache=; platformStoreGroupKey=; isOpenNotification=1; inited=true; console_mode=0; uuid=8454e52abd5eb0de; store_id=85062; vender_id=86; storeId=85062; venderId=86; businessCode=1; appMode=online; storeGroupKey=43c9247c2299f05b6551528fbbf23e89@MS04NTA2Mi04Ng; appVersion=4.9.0; platform=ANDROID; dmall-locale=zh_CN; token=2ed21711-aea4-40c1-a78c-0814aec156d9; ticketName=95E192A474E8E2D41F7B894BE1F58475434F1EBB128F732054FE374078B072EB075399E79F5C7056CB5B90EF2AAA089C5C6AA0C2B8EE37B20135E6B112F9713BDA837A00CA3DA2CE253200765947770CD797C55494673E5D78A80B88FB0884ABDB58FC54C13E2D6D15480DC94300A19DE6C2800E05D54AA08CF5C00746765F7E; userId=206579381; lat=30.706703; lng=111.302199; addr=%E6%B9%96%E5%8C%97%E7%9C%81%E5%AE%9C%E6%98%8C%E5%B8%82%E8%A5%BF%E9%99%B5%E5%8C%BA%E8%A5%BF%E9%99%B5%E8%A1%97%E9%81%93%E5%8A%9E%E4%BA%8B%E5%A4%84; community=%E8%A5%BF%E9%99%B5%E8%A1%97%E9%81%93%E5%8A%9E%E4%BA%8B%E5%A4%84; areaId=420502; session_id=49dd60c7c9f84cebbdddd4ca4b4486f3; env=app; first_session_time=1610966352620; session_count=2; recommend=1; dmTenantId=2"
	tmpToken      = "2ed21711-aea4-40c1-a78c-0814aec156d9"
	tmpTicketName = "95E192A474E8E2D41F7B894BE1F58475434F1EBB128F732054FE374078B072EB075399E79F5C7056CB5B90EF2AAA089C5C6AA0C2B8EE37B20135E6B112F9713BDA837A00CA3DA2CE253200765947770CD797C55494673E5D78A80B88FB0884ABDB58FC54C13E2D6D15480DC94300A19DE6C2800E05D54AA08CF5C00746765F7E"
	tmpUUID       = "8454e52abd5eb0de"
)

func TmpJumpQueue() bool {
	link := "https://presale.dmall.com/maotai/jumpQueue"
	post := url.Values{}

	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"%s","skuId":"1003057062"}`, tmpUUID))
	post.Set("token", tmpToken)
	post.Set("ticketName", tmpTicketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	cookie := fmt.Sprintf(tempCookie, tmpToken, tmpTicketName)
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

	post.Set("param", fmt.Sprintf(`{"vendorId":"86","deviceId":"%s","skuId":"1003057062"}`, tmpUUID))
	post.Set("token", tmpToken)
	post.Set("ticketName", tmpTicketName)

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(post.Encode()))
	if err != nil {
		wlog.Error("创建请求失败:", err)
		return false
	}

	cookie := fmt.Sprintf(tempCookie, tmpToken, tmpTicketName)
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Host", "presale.dmall.com")

	req.Header.Add("Origin", "https://static.dmall.com")
	req.Header.Add("X-Requested-With", "com.dm.metro")
	req.Header.Add("Accept", "application/json, text/plain, */*")

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
