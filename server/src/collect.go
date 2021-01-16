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
