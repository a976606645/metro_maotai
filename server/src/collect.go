package main

import (
	"mdl/wlog"
	"sync"
)

func CollectStock() {
	stores := make([]*TableStockRecord, 0, len(StoreList))
	wg := sync.WaitGroup{}
	for _, item := range StoreList {
		wg.Add(1)
		go func(s StoreInfo) {
			storeInfo, err := queryStock(s.StoreID)
			if err != nil {
				wlog.Error("获取库存失败:", s.StoreID, err)
				return
			}
			stores = append(stores, &TableStockRecord{
				StoreID:   s.StoreID,
				StoreName: s.StoreName,
				Stock:     storeInfo.Data.WareInfo.Stock,
			})
			wg.Done()

		}(item)
	}

	wg.Wait()

	DB.Create(stores)
}

func Start() {
	isRun = true

	users := make([]*TableUser, 0)
	result := DB.Model(&TableUser{}).Where("qualification = true and remain_count > 0").Find(&users)
	if result.Error != nil {
		wlog.Error("数据库读取用户失败:", result.Error)
	}

	for _, user := range users {
		go buy(user)
	}
}

func Stop() {
	isRun = false
	wlog.Info("今日抢购结束")
}

func UpdateUser() {
	users := make([]*TableUser, 0)
	result := DB.Find(users)
	if result.Error != nil {
		wlog.Error("获取用户数据失败:", result.Error)
		return
	}

	for _, user := range users {
		if qinfo, err := qualification(user); err != nil {
			wlog.Errorf("[%s]获取更新用户抢购资格数据失败:%s", user.Phone, err)
			continue
		} else {
			user.Qualification = qinfo.Data.HasQualification
			user.RemainCount = qinfo.Data.RemainCount
			result := DB.Model(&TableUser{}).Where(&TableUser{Phone: user.Phone}).Updates(&TableUser{RemainCount: qinfo.Data.RemainCount, Qualification: qinfo.Data.HasQualification})
			if result.Error != nil {
				wlog.Errorf("[%s]数据库更新用户抢购资格数据失败:%s", user.Phone, err)
			}
		}
	}
}
